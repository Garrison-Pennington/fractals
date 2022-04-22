from abc import ABC, abstractmethod
import PIL.Image as Image

import numpy as np
import matplotlib.pyplot as plt


class Fractal(ABC):
    @abstractmethod
    def __init__(self):
        raise NotImplementedError()

    @abstractmethod
    def step(self, fractal):
        raise NotImplementedError()

    def exit_condition(self, fractal):
        return False

    def n_steps(self, n):
        for i in range(n):
            if self.exit_condition(self.fractal):
                break
            self.fractal = self.step(self.fractal)
        return self.fractal, i+1

    def save(self, img, fname):
        im = Image.fromarray(img * 255).convert("L")
        im.save(fname)

    @abstractmethod
    def render(self, image_shape, steps):
        raise NotImplementedError()

    def score(self, fractal, step):
        raise NotImplementedError()
