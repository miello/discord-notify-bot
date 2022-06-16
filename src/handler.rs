use serenity::framework::standard::macros::group;
use crate::commands::{math::*, misc::*};

#[group]
#[commands(multiply)]
pub struct General;

#[group]
#[commands(ping, latency, test_embed)]
pub struct Misc;

#[group]
#[prefix = "mcv"]
#[commands(test_mcv)]
pub struct McvNotify;