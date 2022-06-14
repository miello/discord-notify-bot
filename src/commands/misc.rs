use serenity::{framework::standard::{macros::command, Args, CommandResult}, client::{Context, bridge::gateway::ShardId}, model::{channel::Message}, utils::MessageBuilder};
use crate::{ShardManagerContainer, utils::{http::get_raw_http_response, env::get_env}};
use scraper::{Html, Selector};
use chrono::{prelude::*, Duration};

#[derive(Debug)]
struct Course {
    course_no: String,
    course_title: String,
    course_year: String,
    course_semester: String,
    course_href: String
}

impl Course {
    pub fn get_title(&self) -> String {
        format!("{} ({}/{})", self.course_no, self.course_year, self.course_semester)
    } 

    pub fn get_description(&self) -> String {
        format!("{}", self.course_title)
    }
}

#[command]
pub async fn ping(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let find_channel = match msg.channel_id.to_channel(&ctx.http).await {
        Ok(channel) => Some(channel),
        Err(why) => {
            println!("Error getting channel: {:?}", why);
            None
        },
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
            msg.reply(ctx, "There was a problem getting the shard manager").await?;

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
        },
    };

    match runner.latency {
        Some(latency) => {
            msg.reply(ctx, &format!("The shard latency is {:?}", latency)).await?;
        },
        None => {
            msg.reply(ctx, "Error when get latency from shard").await?;
        }
    }


    Ok(())
}

#[command]
pub async fn test_embed(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    msg.channel_id.send_message(&ctx.http, |m| {
        m.embed(|e| {
            e.title("Daily MCV");
            e.description("This is a description of the embed!");
            e.color(0x018ada);
            e.url("https://google.com");
            e.thumbnail(format!("{}sites/all/modules/courseville/files/logo/cv-logo.png", get_env("BASE_URL")));
            e.field("test", 123, false);
            e.field("test", 12312321, false);
            e.field("test", "asdas", false);

            e
        });

        

        m
    }).await?;

    Ok(())
}

#[command]
pub async fn test_mcv(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("BASE_URL");
    let text = get_raw_http_response(&format!("{}?q=courseville", base_url)).await.unwrap();

    // Tokio try to make it thread-safe but Html does not support 'Send' impl
    let get_all_course = || {
        let selector = Selector::parse("*[course_no]").unwrap();
        Html::parse_document(&text).select(&selector).map(|f| {
            let value = f.value();
            let get_key = |key: &str| {
                value.attr(key).unwrap_or(&String::new()).to_string()
            };
            Course {
                course_no: get_key("course_no"),
                course_title: get_key("title"),
                course_href: format!("{}{}", base_url, get_key("href")),
                course_semester: get_key("semester"),
                course_year: get_key("year"),
            }
        }).collect::<Vec<Course>>()
    };

    let all_course = get_all_course();

    msg.channel_id.send_message(&ctx.http, |m| {
        m.embed(move |e| {
            e.title("MCV Notify");
            e.thumbnail(format!("{}sites/all/modules/courseville/files/logo/cv-logo.png", base_url));

            all_course.iter().for_each(|course| {
                e.field(course.get_title(), course.get_description(), true);
            });

            e.footer(|footer| {
                let current_time =  Utc::now() + Duration::hours(7);
                footer.text(format!("Update at {}", current_time.to_rfc3339()))
            })
        })
    }).await?;

    Ok(())
}