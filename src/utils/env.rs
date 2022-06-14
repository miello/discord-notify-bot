use std::env;

pub fn get_env(key: &str) -> String {
    env::var(key).unwrap_or(String::new())
}

pub fn get_env_required(key: &str) -> String {
    env::var(key).expect(&format!("Expected {}", key))
}