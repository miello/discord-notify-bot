use crate::{
    utils::{env::get_env},
    ShardManagerContainer,
};
use serenity::{
    client::{bridge::gateway::ShardId, Context},
    framework::standard::{macros::command, Args, CommandResult},
    model::channel::Message, utils::MessageBuilder,
};

#[command]
pub async fn ping(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let find_channel = match msg.channel_id.to_channel(&ctx.http).await {
        Ok(channel) => Some(channel),
        Err(why) => {
            println!("Error getting channel: {:?}", why);
            None
        }
    };

    let channel = find_channel.unwrap();

    let response = MessageBuilder::new()
        .push("User ")
        .push_bold_safe(&msg.author.name)
        .push(" used the 'ping' command in the ")
        .mention(&channel)
        .push(" channel")
        .build();

    if let Err(why) = msg.channel_id.say(&ctx.http, &response).await {
        println!("Error sending message: {:?}", why);
    };

    Ok(())
}

#[command]
pub async fn latency(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let data = ctx.data.read().await;

    let shard_manager = match data.get::<ShardManagerContainer>() {
        Some(v) => v,
        None => {
            msg.reply(ctx, "There was a problem getting the shard manager")
                .await?;

            return Ok(());
        }
    };

    let manager = shard_manager.lock().await;
    let runners = manager.runners.lock().await;

    let runner = match runners.get(&ShardId(ctx.shard_id)) {
        Some(runner) => runner,
        None => {
            msg.reply(ctx, "No shard found").await?;

            return Ok(());
        }
    };

    match runner.latency {
        Some(latency) => {
            msg.reply(ctx, &format!("The shard latency is {:?}", latency))
                .await?;
        }
        None => {
            msg.reply(ctx, "Error when get latency from shard").await?;
        }
    }

    Ok(())
}

#[command]
pub async fn test_embed(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    msg.channel_id
        .send_message(&ctx.http, |m| {
            m.embed(|e| {
                e.title("Daily MCV");
                e.description("This is a description of the embed!");
                e.color(0x018ada);
                e.url("https://google.com");
                e.thumbnail(format!(
                    "{}sites/all/modules/courseville/files/logo/cv-logo.png",
                    get_env("MCV_BASE_URL")
                ));
                e.field("test", 123, false);
                e.field("test", 12312321, false);
                e.field("test", "asdas", false);

                e
            });

            m
        })
        .await?;

    Ok(())
}

