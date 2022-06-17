use crate::{
    utils::{env::get_env, http::get_raw_http_response_mcv},
    ShardManagerContainer,
};
use chrono::{prelude::*, Duration};
use scraper::{ElementRef, Html, Selector};
use serenity::{
    client::{bridge::gateway::ShardId, Context},
    framework::standard::{macros::command, Args, CommandResult},
    model::channel::Message,
    utils::{EmbedMessageBuilding, MessageBuilder},
};
use tracing::log::{log, info};

#[derive(Debug)]
struct Course {
    course_no: String,
    course_title: String,
    course_year: i16,
    course_semester: i16,
    course_href: String,
}

impl Course {
    pub fn get_title(&self) -> String {
        format!(
            "{} ({}/{})",
            self.course_no, self.course_year, self.course_semester
        )
    }

    pub fn get_description(&self) -> String {
        format!("{}", self.course_title)
    }
}

struct Announcement {
    title: String,
    date: String,
    href: String,
}

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

#[command("all_course")]
pub async fn test_mcv(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("MCV_BASE_URL");
    let text = get_raw_http_response_mcv(&format!("{}?q=courseville", base_url))
        .await
        .unwrap();
    if text.contains("Please login with either of the following choices.") {
        msg.channel_id
            .send_message(&ctx.http, |f| {
                f.content("Please contact bot creator to login MCV")
            })
            .await?;
        return Ok(());
    }

    // Tokio try to make it thread-safe but Html does not support 'Send' impl
    let get_all_course = || {
        let selector = Selector::parse("*[course_no]").unwrap();
        Html::parse_document(&text)
            .select(&selector)
            .map(|f| {
                let value = f.value();
                let get_key = |key: &str| value.attr(key).unwrap_or(&String::new()).to_string();
                Course {
                    course_no: get_key("course_no"),
                    course_title: get_key("title"),
                    course_href: format!("{}{}", base_url, get_key("href")),
                    course_semester: get_key("semester").parse::<i16>().unwrap(),
                    course_year: get_key("year").parse::<i16>().unwrap(),
                }
            })
            .collect::<Vec<Course>>()
    };

    let all_course = get_all_course();
    let filter_course: Vec<&Course> = all_course
        .iter()
        .filter(|course| {
            let msg_list: Vec<&str> = msg.content.split(" ").collect();
            let msg_len = msg_list.len();
            if msg_len <= 2 {
                return true;
            }

            let course_year: i16 = msg_list[2].parse().unwrap_or(-1);
            if msg_len <= 3 {
                if course_year == course.course_year {
                    return true;
                }
                return false;
            }

            let course_semester: i16 = msg_list[3].parse().unwrap_or(-1);
            if msg_len <= 4 {
                if course_year == course.course_year && course_semester == course.course_semester {
                    return true;
                }
                return false;
            }
            false
        })
        .collect();

    msg.channel_id
        .send_message(&ctx.http, |m| {
            m.embed(move |e| {
                e.title("MCV Notify");
                e.thumbnail(format!(
                    "{}sites/all/modules/courseville/files/logo/cv-logo.png",
                    base_url
                ));

                filter_course.iter().for_each(|course| {
                    e.field(course.get_title(), course.get_description(), true);
                });

                e.footer(|footer| {
                    let current_time = Utc::now() + Duration::hours(7);
                    footer.text(format!("Update at {}", current_time.to_rfc3339()))
                })
            })
        })
        .await?;

    Ok(())
}

#[command("announce")]
pub async fn test_get_announcement(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("MCV_BASE_URL");
    let split_msg: Vec<&str> = msg.content.split(" ").collect();

    if split_msg.len() != 3 {
        msg.channel_id
            .send_message(&ctx.http, |c| c.content("Expected course id"))
            .await?;
        return Ok(());
    }
    let course_id = split_msg[2];

    // let text = get_raw_http_response_mcv("http://localhost:5500/src")
    //     .await
    //     .unwrap();

    let text = get_raw_http_response_mcv(&format!("{}?q=courseville/course/{}", base_url, course_id))
        .await
        .unwrap();
        
    if text.contains("Please login with either of the following choices.") {
        msg.channel_id
            .send_message(&ctx.http, |f| {
                f.content("Please contact bot creator to login MCV")
            })
            .await?;
        return Ok(());
    }

    if text.contains("It looks like you are not a member of this course yet.") {
        msg.channel_id
            .send_message(&ctx.http, |f| {
                f.content("Bot owner does not enroll this course")
            })
            .await?;
        return Ok(());
    }

    let get_all_annoucement = |html: &str| {
        
        let selector = Selector::parse("table[title='Course announcements'] > tbody > tr").unwrap();
        let result = Html::parse_document(&html);
        let title_el = result.select(&selector).collect::<Vec<_>>();

        title_el
            .iter()
            .map(|tr| {
                let selector_td = Selector::parse("td").unwrap();
                let mut tr_iter = tr.select(&selector_td);

                let date_root = tr_iter.next().unwrap().first_child().unwrap();
                let desc_root = tr_iter.next().unwrap().first_child().unwrap();

                // Date string
                let date = ElementRef::wrap(date_root).unwrap().inner_html();
                let title = ElementRef::wrap(desc_root).unwrap().inner_html();
                let href = format!("{}{}", base_url, &desc_root.value().as_element().unwrap().attr("href").unwrap_or("").to_string());

                Announcement { date, href, title }
            })
            .collect::<Vec<Announcement>>()
    };

    let get_course_title = |html: &str| {
        let selector = Selector::parse(".courseville-course-title").unwrap();
        let result = Html::parse_document(&html);
        let title_div = ElementRef::wrap(result.select(&selector).next().unwrap().first_child().unwrap()).unwrap();

        title_div.inner_html()
    };

    let course_title = get_course_title(&text);
    let all_announcement = get_all_annoucement(&text);

    msg.channel_id
        .send_message(&ctx.http, |m| {
            m.embed(|e| {
                let title = format!("{} Announcement", course_title);
                e.title(&title);

                all_announcement.iter().for_each(|announce| {
                    let mut builder = MessageBuilder::new();

                    let desc = builder.push_named_link_safe(&announce.title, &announce.href);
                    e.field(&announce.date, &desc, false);
                });

                e
            });

            m
        })
        .await?;

    Ok(())
}
