use wasm_bindgen::prelude::*;
use photon_rs::{native::open_image_from_bytes, transform, conv, monochrome};
use image::{codecs::png::PngEncoder, ColorType, ImageEncoder};
use std::io::Cursor;

#[wasm_bindgen(start)]
pub fn init() {
    #[cfg(debug_assertions)]
    console_error_panic_hook::set_once();
}

#[wasm_bindgen]
pub fn process_image(
    input_bytes: &[u8],
    target_width: u32,
    to_gray: bool,
) -> Result<Vec<u8>, JsValue> {
    let mut img = decode_image(input_bytes)?;

    resize_image(&mut img, target_width);

    if to_gray {
        apply_grayscale_and_sharpen(&mut img);
    }

    encode_png(&img)
}

fn decode_image(bytes: &[u8]) -> Result<photon_rs::PhotonImage, JsValue> {
    open_image_from_bytes(bytes).map_err(|_| JsValue::from_str("이미지 디코딩 실패"))
}

fn resize_image(img: &mut photon_rs::PhotonImage, target_width: u32) {
    let (w, h) = (img.get_width(), img.get_height());
    if target_width > 0 && target_width < w {
        let new_h = ((h as f32) * (target_width as f32 / w as f32)).round() as u32;
        transform::resize(img, target_width, new_h, transform::SamplingFilter::CatmullRom);
    }
}

fn apply_grayscale_and_sharpen(img: &mut photon_rs::PhotonImage) {
    monochrome::grayscale_shades(img, 4);
    conv::sharpen(img);
}

fn encode_png(img: &photon_rs::PhotonImage) -> Result<Vec<u8>, JsValue> {
    let (w, h) = (img.get_width(), img.get_height());
    let rgba = img.get_raw_pixels();

    let mut out = Vec::new();
    let mut cursor = Cursor::new(&mut out);
    let encoder = PngEncoder::new(&mut cursor);

    encoder
        .write_image(&rgba, w, h, ColorType::Rgba8)
        .map_err(|e| JsValue::from_str(&format!("PNG 인코딩 실패: {e}")))?;

    Ok(out)
}