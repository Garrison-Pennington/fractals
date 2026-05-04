//! wasm-pack integration tests — run with `wasm-pack test --headless --chrome`.

#![cfg(target_arch = "wasm32")]

use fractals_wasm::*;
use wasm_bindgen::JsValue;
use wasm_bindgen_test::*;

// Tests run in both Node (wasm-pack test --node) and headless browsers
// (wasm-pack test --headless --chrome). The js-sys helpers used below are
// available in both environments.

#[wasm_bindgen_test]
fn renders_mandelbrot_region() {
    // 100x100 region centered at (-0.5, 0) — matches acceptance criterion (a).
    let buf = render_fractal(
        "mandelbrot", 0.0, 0.0, -1.5, 0.5, -1.0, 1.0, 100, 100, 200,
    );
    assert_eq!(buf.len(), 100 * 100 * 2);
    let first = buf[0];
    let any_different = buf.iter().step_by(2).any(|v| (*v - first).abs() > 1e-6);
    assert!(any_different, "expected non-uniform pixel data");
}

#[wasm_bindgen_test]
fn interior_point_returns_max_iterations() {
    // Criterion (b): point (0, 0) should stay inside for max_iter steps.
    let info = compute_point("mandelbrot", 0.0, 0.0, 0.0, 0.0, 500);
    let iters = js_sys::Reflect::get(&info, &JsValue::from_str("iterations"))
        .unwrap()
        .as_f64()
        .unwrap();
    let escaped = js_sys::Reflect::get(&info, &JsValue::from_str("escaped"))
        .unwrap()
        .as_bool()
        .unwrap();
    assert_eq!(iters as u32, 500);
    assert!(!escaped);
}

#[wasm_bindgen_test]
fn exterior_point_escapes_quickly() {
    // Criterion (c): point (2, 2) escapes within a few iterations.
    let info = compute_point("mandelbrot", 0.0, 0.0, 2.0, 2.0, 500);
    let iters = js_sys::Reflect::get(&info, &JsValue::from_str("iterations"))
        .unwrap()
        .as_f64()
        .unwrap() as u32;
    let escaped = js_sys::Reflect::get(&info, &JsValue::from_str("escaped"))
        .unwrap()
        .as_bool()
        .unwrap();
    assert!(escaped);
    assert!(iters < 5);
}

#[wasm_bindgen_test]
fn music_sequence_has_expected_count_and_values() {
    // Criterion (d): path-based music sequence returns the requested count
    // with valid pitch/duration/velocity values.
    let cfg = js_sys::Object::new();
    let set = |k: &str, v: JsValue| {
        js_sys::Reflect::set(&cfg, &JsValue::from_str(k), &v).unwrap();
    };
    set("kind", JsValue::from_str("mandelbrot"));
    set("c_re", JsValue::from_f64(0.0));
    set("c_im", JsValue::from_f64(0.0));
    set("scale", JsValue::from_str("major"));
    set("root_midi", JsValue::from_f64(60.0));
    set("sequence_type", JsValue::from_str("path"));
    set("length", JsValue::from_f64(16.0));
    set("max_iter", JsValue::from_f64(200.0));
    set("x_min", JsValue::from_f64(-2.0));
    set("x_max", JsValue::from_f64(1.0));
    set("y_min", JsValue::from_f64(-0.3));
    set("y_max", JsValue::from_f64(0.3));

    let notes = generate_music(cfg.into());
    let arr = js_sys::Array::from(&notes);
    assert_eq!(arr.length(), 16);
    for i in 0..arr.length() {
        let n = arr.get(i);
        let pitch = js_sys::Reflect::get(&n, &JsValue::from_str("pitch"))
            .unwrap()
            .as_f64()
            .unwrap() as u32;
        let duration = js_sys::Reflect::get(&n, &JsValue::from_str("duration"))
            .unwrap()
            .as_f64()
            .unwrap();
        let velocity = js_sys::Reflect::get(&n, &JsValue::from_str("velocity"))
            .unwrap()
            .as_f64()
            .unwrap() as u32;
        assert!(pitch < 128);
        assert!(duration > 0.0);
        assert!(velocity > 0 && velocity < 128);
    }
}

#[wasm_bindgen_test]
fn progressive_matches_full_resolution_dimensions() {
    let buf =
        render_fractal_progressive("mandelbrot", 0.0, 0.0, -2.0, 1.0, -1.0, 1.0, 64, 64, 200, 4);
    assert_eq!(buf.len(), 64 * 64 * 2);
}

#[wasm_bindgen_test]
fn default_bounds_are_sensible() {
    let b = default_bounds("mandelbrot");
    let x_min = js_sys::Reflect::get(&b, &JsValue::from_str("x_min"))
        .unwrap()
        .as_f64()
        .unwrap();
    let x_max = js_sys::Reflect::get(&b, &JsValue::from_str("x_max"))
        .unwrap()
        .as_f64()
        .unwrap();
    assert!(x_min < 0.0 && x_max > 0.0);
    assert!(x_max - x_min > 1.0);
}
