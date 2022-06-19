use scraper::{Selector, Html, ElementRef};

#[derive(Debug)]
pub struct Course {
    pub course_id: String,
    pub course_no: String,
    pub course_title: String,
    pub course_year: i16,
    pub course_semester: i16,
    pub course_href: String,
}

impl Course {
    pub fn get_title(&self) -> String {
        format!(
            "[{}] {} ({}/{})",
            self.course_id, self.course_no, self.course_year, self.course_semester
        )
    }

    pub fn get_description(&self) -> String {
        format!("{}", self.course_title)
    }
}

pub struct Announcement {
    pub title: String,
    pub date: String,
    pub href: String,
}

pub struct Assignment {
    pub title: String,
    pub due_date: String,
    pub href: String,
}

impl Assignment {
    pub fn get_title(&self) -> String {
        format!(
            "{}",
            self.title
        )
    }

    pub fn get_description(&self) -> String {
        format!("{} ", self.due_date)
    }
}

pub fn get_all_assignment(html: &str, base_url: &str, limit: u32) -> Vec<Assignment> {
    let selector = Selector::parse("table[title='Assignment list'] > tbody > tr").unwrap();
    let result = Html::parse_document(&html);
    let title_el = result.select(&selector).collect::<Vec<_>>();

    title_el
        .iter()
        .map(|tr| {
            let selector_td = Selector::parse("td").unwrap();
            let select_a = Selector::parse("a").unwrap();
            let mut td_iter = tr.select(&selector_td);

            // Skip first column
            td_iter.next();

            let title_col = td_iter.next().unwrap().select(&select_a).next().unwrap();
            let title_el = title_col.value();

            let href = format!("{}{}", base_url, title_el.attr("href").unwrap());
            let title = title_col.inner_html();

            // Skip out date column
            td_iter.next();

            // Get due date column
            let selector_sr_only = Selector::parse(".sr-only").unwrap();
            let due_date = td_iter
                .next()
                .unwrap()
                .select(&selector_sr_only)
                .next()
                .unwrap()
                .inner_html();

            Assignment {
                title,
                due_date,
                href,
            }
        })
        .collect()
}

pub fn get_course_title(html: &str) -> String {
    let selector = Selector::parse(".courseville-course-title").unwrap();
    let result = Html::parse_document(&html);
    let title_div = ElementRef::wrap(
        result
            .select(&selector)
            .next()
            .unwrap()
            .first_child()
            .unwrap(),
    )
    .unwrap();

    title_div.inner_html()
}

pub fn get_all_annoucement(html: &str, base_url: &str) -> Vec<Announcement> {
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
            let href = format!(
                "{}{}",
                base_url,
                &desc_root
                    .value()
                    .as_element()
                    .unwrap()
                    .attr("href")
                    .unwrap_or("")
                    .to_string()
            );

            Announcement { date, href, title }
        })
        .collect()
}

pub fn get_all_course(html: &str, base_url: &str) -> Vec<Course> {
    let selector = Selector::parse("*[course_no]").unwrap();
    Html::parse_document(&html)
        .select(&selector)
        .map(|f| {
            let value = f.value();
            let get_key = |key: &str| value.attr(key).unwrap_or(&String::new()).to_string();
            Course {
                course_id: get_key("cv_cid"),
                course_no: get_key("course_no"),
                course_title: get_key("title"),
                course_href: format!("{}{}", base_url, get_key("href")),
                course_semester: get_key("semester").parse::<i16>().unwrap(),
                course_year: get_key("year").parse::<i16>().unwrap(),
            }
        })
        .collect()
}