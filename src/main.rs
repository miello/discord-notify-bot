mod commands;
mod handler;
mod utils;

use std::{collections::HashSet, env, sync::Arc};

use crate::handler::*;
use dotenv::dotenv;
use serenity::{
    async_trait,
    client::{bridge::gateway::ShardManager, Context, EventHandler},
    framework::StandardFramework,
    http::Http,
    model::{channel::Message, event::ResumedEvent, gateway::Ready},
    prelude::{GatewayIntents, Mutex, TypeMapKey},
    Client,
};
use tracing::{error, info};
use utils::env::get_env_required;

pub struct ShardManagerContainer;

impl TypeMapKey for ShardManagerContainer {
    type Value = Arc<Mutex<ShardManager>>;
}

struct Handler;

#[async_trait]
impl EventHandler for Handler {
    async fn ready(&self, _: Context, ready: Ready) {
        info!("Connected as {}", ready.user.name);
    }

    async fn resume(&self, _: Context, _: ResumedEvent) {
        info!("Resumed");
    }

    async fn message(&self, _ctx: Context, _new_message: Message) {
        if _new_message.author.bot {
            return;
        }
        info!("Got the message: {}", _new_message.content);
    }
}

#[tokio::main]
async fn main() {
    dotenv().expect("Failed to load .env file");

    tracing_subscriber::fmt::init();

    let token = get_env_required("DISCORD_TOKEN");

    let http = Http::new(&token);
    let (owners, _bot_id) = match http.get_current_application_info().await {
        Ok(info) => {
            let mut owners = HashSet::new();
            owners.insert(info.owner.id);

            (owners, info.id)
        }
        Err(why) => panic!("Could not access application info: {:?}", why),
    };

    let framework = StandardFramework::new()
        .configure(|c| c.owners(owners).prefix('~'))
        .group(&GENERAL_GROUP)
        .group(&MISC_GROUP);

    let intents = GatewayIntents::GUILD_MESSAGES
        | GatewayIntents::DIRECT_MESSAGES
        | GatewayIntents::MESSAGE_CONTENT;

    let mut client = Client::builder(&token, intents)
        .framework(framework)
        .event_handler(Handler)
        .await
        .expect("Error on create client");

    {
        let mut data = client.data.write().await;
        data.insert::<ShardManagerContainer>(client.shard_manager.clone());
    }

    let shard_manager = client.shard_manager.clone();

    // Gracefully shutdown
    tokio::spawn(async move {
        tokio::signal::ctrl_c()
            .await
            .expect("Could not register ctrl+c handler");
        info!("Gracefully shutdown");
        shard_manager.lock().await.shutdown_all().await;
    });

    if let Err(why) = client.start().await {
        error!("Client Error: {:?}", why)
    }
}
