use serenity::{framework::standard::{macros::command, Args, CommandResult}, client::Context, model::channel::Message};

#[command]
pub async fn multiply(ctx: &Context, msg: &Message, mut args: Args) -> CommandResult {
    let one = args.single::<i64>()?;
    let two = args.single::<i64>()?;

    let product = one * two;
    
    msg.channel_id.say(&ctx.http, product).await?;

    Ok(())
}