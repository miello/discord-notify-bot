use crate::utils::helper::*;
use crate::utils::{env::get_env, http::get_raw_http_response_mcv};
use chrono::{prelude::*, Duration};
use serenity::framework::standard::macros::command;
use serenity::{
    client::Context,
    framework::standard::{Args, CommandResult},
    model::channel::Message,
    utils::{EmbedMessageBuilding, MessageBuilder},
};

fn get_error_message(text: &str) -> Option<String> {
    if text.contains("Please login with either of the following choices.") {
        return Some(String::from("Please contact bot creator to login MCV"));
    }

    if text.contains("It looks like you are not a member of this course yet.") {
        return Some(String::from("Bot owner does not enroll this course"));
    }

    None
}

#[command("all_course")]
pub async fn get_courses(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("MCV_BASE_URL");
    let text = get_raw_http_response_mcv(&format!("{}?q=courseville", base_url))
        .await
        .unwrap();

    if let Some(val) = get_error_message(&text) {
        msg.channel_id
            .send_message(&ctx.http, |f| f.content(val))
            .await?;
    }

    let all_course = get_all_course(&text, &base_url);
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
                e.url(format!("{}?q=courseville", base_url));
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
pub async fn get_announcement(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("MCV_BASE_URL");
    let split_msg: Vec<&str> = msg.content.split(" ").collect();

    if split_msg.len() != 3 {
        msg.channel_id
            .send_message(&ctx.http, |c| c.content("Expected course id"))
            .await?;
        return Ok(());
    }

    let course_id = split_msg[2];

    let text =
        get_raw_http_response_mcv(&format!("{}?q=courseville/course/{}", base_url, course_id))
            .await
            .unwrap();

    if let Some(val) = get_error_message(&text) {
        msg.channel_id
            .send_message(&ctx.http, |f| f.content(val))
            .await?;
    }

    let all_announcement = get_all_annoucement(&text, &base_url);
    let course_title = get_course_title(&text);

    msg.channel_id
        .send_message(&ctx.http, |m| {
            m.embed(|e| {
                let title = format!("{} Announcement", course_title);

                e.title(&title);
                e.url(format!("{}?q=courseville/course/{}", base_url, course_id));
                e.thumbnail(format!(
                    "{}sites/all/modules/courseville/files/logo/cv-logo.png",
                    base_url
                ));

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

#[command("assign")]
pub async fn get_assignment(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("MCV_BASE_URL");
    let split_msg: Vec<&str> = msg.content.split(" ").collect();

    if split_msg.len() != 3 {
        msg.channel_id
            .send_message(&ctx.http, |c| c.content("Expected course id"))
            .await?;
        return Ok(());
    }
    let course_id = split_msg[2];

    let text = get_raw_http_response_mcv(&format!(
        "{}?q=courseville/course/{}/assignment",
        base_url, course_id
    ))
    .await
    .unwrap();

    if let Some(val) = get_error_message(&text) {
        msg.channel_id
        .send_message(&ctx.http, |f| {
            f.content(val)
        })
        .await?;
    }

    let all_assignment = get_all_assignment(&text, &base_url, 5);
    let course_title = get_course_title(&text);

    msg.channel_id
        .send_message(&ctx.http, |m| {
            m.embed(|e| {
                let title = format!("{} Assignment", course_title);

                e.title(&title);
                e.url(format!(
                    "{}?q=courseville/course/{}/assignment",
                    base_url, course_id
                ));
                e.thumbnail(format!(
                    "{}sites/all/modules/courseville/files/logo/cv-logo.png",
                    base_url
                ));

                all_assignment.iter().for_each(|assignment| {
                    let mut desc = MessageBuilder::new();

                    desc.push_named_link_safe(assignment.get_description(), &assignment.href);
                    
                    e.field(&assignment.get_title(), &desc, false);
                });

                e
            });

            m
        })
        .await?;

    Ok(())
}

#[command("file")]
pub async fn get_materials(ctx: &Context, msg: &Message, _: Args) -> CommandResult {
    let base_url: String = get_env("MCV_BASE_URL");
    let split_msg: Vec<&str> = msg.content.split(" ").collect();

    if split_msg.len() != 3 {
        msg.channel_id
            .send_message(&ctx.http, |c| c.content("Expected course id"))
            .await?;
        return Ok(());
    }
    let course_id = split_msg[2];

    let text = get_raw_http_response_mcv(&format!(
        "{}?q=courseville/course/{}",
        base_url, course_id
    ))
        .await
        .unwrap();

    if let Some(val) = get_error_message(&text) {
        msg.channel_id
        .send_message(&ctx.http, |f| {
            f.content(val)
        })
        .await?;

        return Ok(());
    }

    let materials = get_all_material(&text, &base_url);
    let course_title = get_course_title(&text);
    
    let res = msg.channel_id.send_message(&ctx.http, |m| {
        m.embed(|e| {
            let title = format!("{} Materials", course_title);
            e.title(&title);
            e.url(&format!(
                "{}?q=courseville/course/{}",
                base_url, course_id
            ));
            e.thumbnail(format!(
                "{}sites/all/modules/courseville/files/logo/cv-logo.png",
                base_url
            ));

            materials.iter().for_each(|folder| {
                folder.embed_message(e, 3)
            });

            e
        })
    }).await;

    Ok(())
}
