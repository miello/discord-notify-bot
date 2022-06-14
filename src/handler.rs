use serenity::framework::standard::macros::group;
use crate::commands::{math::*, misc::*};

#[group]
#[commands(multiply)]
pub struct General;

#[group]
#[commands(ping, latency, test_embed, test_mcv)]
pub struct Misc;