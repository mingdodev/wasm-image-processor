use wasm_bindgen::prelude::*;
use photon_rs::{native::open_image_from_bytes, transform, conv, monochrome};

#[wasm_bindgen(start)]
pub fn init() {
    #[cfg(debug_assertions)]
    console_error_panic_hook::set_once();
}

#[wasm_bindgen]
pub fn process_image(input_bytes: &[u8], target_width: u32, to_gray: bool) -> Vec<u8> {
    let mut img = open_image_from_bytes(input_bytes).expect("decode failed");

    let (w, h) = (img.get_width(), img.get_height());
    if target_width > 0 && target_width < w {
        let new_h = (h as f32 * (target_width as f32 / w as f32)).round() as u32;
        transform::resize(&mut img, target_width, new_h, transform::SamplingFilter::CatmullRom);
    }

    if to_gray {
        monochrome::grayscale_shades(&mut img, 4_u8);
        conv::sharpen(&mut img);
    }

    photon_rs::native::image_to_bytes(img)
}
