use reqwest::header::{COOKIE, HeaderMap, HeaderValue};

const SESSION_KEY: &str = "SESSeb912a58562fbbdf6ad5e9a19524d1c0";
const SESSION_VALUE: &str = "mf1cfn2vhl1f2srjugrv8qvov3";
pub async fn get_raw_http_response(url: &str) -> Result<String, Box<dyn std::error::Error>> {
    let mut request_headers = HeaderMap::new();
    
    let cookie_header = format!("{}={}", SESSION_KEY, SESSION_VALUE);
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
