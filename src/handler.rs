use serenity::framework::standard::macros::group;
use crate::commands::{math::*, misc::*, mcv_scraper::*};

#[group]
#[commands(multiply)]
pub struct General;

#[group]
#[commands(ping, latency, test_embed)]
pub struct Misc;

#[group]
#[prefix = "mcv"]
#[commands(get_courses, get_announcement, get_assignment)]
pub struct McvNotify;