use reqwest::header::{COOKIE, HeaderMap, HeaderValue};

use crate::utils::env::get_env;

pub async fn get_raw_http_response_mcv(url: &str) -> Result<String, Box<dyn std::error::Error>> {
    let mut request_headers = HeaderMap::new();
    let session_key: String = get_env("MCV_SESSION_KEY");
    let session_value: String = get_env("MCV_SESSION_VALUE");

    let cookie_header = format!("{}={}", session_key, session_value);
    request_headers.insert(COOKIE, HeaderValue::from_str(&cookie_header).unwrap());
    let client = reqwest::ClientBuilder::new()
            .default_headers(request_headers)
            .cookie_store(true)
            .build()
            .unwrap();

    let resp = client.get(url).send().await?;
    let text = resp.text().await?;

    Ok(text)
}
