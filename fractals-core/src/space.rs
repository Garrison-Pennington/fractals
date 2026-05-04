//! Coordinate space transforms between pixel and complex planes.

use num_complex::Complex64;

/// Rectangular region of the complex plane.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct Bounds {
    pub x_min: f64,
    pub x_max: f64,
    pub y_min: f64,
    pub y_max: f64,
}

impl Bounds {
    pub fn new(x_min: f64, x_max: f64, y_min: f64, y_max: f64) -> Self {
        Self { x_min, x_max, y_min, y_max }
    }

    pub fn width(&self) -> f64 {
        self.x_max - self.x_min
    }

    pub fn height(&self) -> f64 {
        self.y_max - self.y_min
    }

    pub fn center(&self) -> Complex64 {
        Complex64::new((self.x_min + self.x_max) * 0.5, (self.y_min + self.y_max) * 0.5)
    }

    /// Zoom in/out by `scale` centered on `focal`. `scale > 1` zooms in.
    pub fn zoom(&self, focal: Complex64, scale: f64) -> Self {
        let w = self.width() / scale;
        let h = self.height() / scale;
        Self {
            x_min: focal.re - w * 0.5,
            x_max: focal.re + w * 0.5,
            y_min: focal.im - h * 0.5,
            y_max: focal.im + h * 0.5,
        }
    }

    /// Convert a pixel coordinate to a complex number, given image resolution.
    /// Pixel (0, 0) is top-left. Y axis of the complex plane increases downward
    /// (i.e. y_min corresponds to the top of the image), matching canvas
    /// conventions so callers don't have to flip rows.
    pub fn pixel_to_complex(&self, px: u32, py: u32, res: Resolution) -> Complex64 {
        let u = (px as f64 + 0.5) / res.width as f64;
        let v = (py as f64 + 0.5) / res.height as f64;
        Complex64::new(self.x_min + u * self.width(), self.y_min + v * self.height())
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Resolution {
    pub width: u32,
    pub height: u32,
}

impl Resolution {
    pub fn new(width: u32, height: u32) -> Self {
        Self { width, height }
    }

    pub fn total_pixels(&self) -> u32 {
        self.width * self.height
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn zoom_in_halves_shape() {
        let b = Bounds::new(-2.0, 2.0, -2.0, 2.0);
        let z = b.zoom(Complex64::new(0.0, 0.0), 2.0);
        assert!((z.width() - 2.0).abs() < 1e-12);
        assert!((z.height() - 2.0).abs() < 1e-12);
        assert!((z.x_min + 1.0).abs() < 1e-12);
    }

    #[test]
    fn pixel_to_complex_corners() {
        let b = Bounds::new(-2.0, 1.0, -1.5, 1.5);
        let r = Resolution::new(100, 100);
        let top_left = b.pixel_to_complex(0, 0, r);
        let bottom_right = b.pixel_to_complex(99, 99, r);
        assert!(top_left.re > -2.0 && top_left.re < -1.97);
        assert!(top_left.im > -1.5 && top_left.im < -1.47);
        assert!(bottom_right.re < 1.0 && bottom_right.re > 0.97);
        assert!(bottom_right.im < 1.5 && bottom_right.im > 1.47);
    }
}
