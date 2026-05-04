// Web Audio synthesizer with an ADSR envelope. Plays a `GeneratedNote[]`
// sequence and fires a callback for each note-on / note-off event so the UI
// can animate a playback cursor.

import type { GeneratedNote, Waveform } from "./types.js";

export interface PlayOptions {
  waveform: Waveform;
  bpm: number;
  onNote?: (index: number, note: GeneratedNote) => void;
  onEnd?: () => void;
}

export interface PlayHandle {
  stop(): void;
  get playing(): boolean;
}

export function midiToFreq(midi: number): number {
  return 440 * Math.pow(2, (midi - 69) / 12);
}

/**
 * Play a sequence using one oscillator per note with an ADSR envelope.
 * Scheduling is done up-front via AudioContext.currentTime; stopping cancels
 * all pending oscillators.
 */
export function playSequence(
  notes: GeneratedNote[],
  opts: PlayOptions
): PlayHandle {
  if (typeof window === "undefined") {
    return { stop() {}, get playing() { return false; } };
  }
  const AudioCtor: typeof AudioContext =
    (window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext);
  const ctx = new AudioCtor();
  const master = ctx.createGain();
  master.gain.value = 0.25;
  master.connect(ctx.destination);

  const secPerBeat = 60 / opts.bpm;
  const now = ctx.currentTime + 0.05;
  const oscillators: OscillatorNode[] = [];
  const timers: number[] = [];
  let stopped = false;

  let cursor = now;
  notes.forEach((note, i) => {
    const dur = note.duration * secPerBeat;
    const attack = Math.min(0.02, dur * 0.2);
    const decay = Math.min(0.05, dur * 0.2);
    const sustainLevel = 0.6;
    const release = Math.min(0.05, dur * 0.2);
    const peak = Math.max(0.05, note.velocity / 127);

    const osc = ctx.createOscillator();
    osc.type = opts.waveform;
    osc.frequency.value = midiToFreq(note.pitch);
    const gain = ctx.createGain();
    gain.gain.setValueAtTime(0, cursor);
    gain.gain.linearRampToValueAtTime(peak, cursor + attack);
    gain.gain.linearRampToValueAtTime(peak * sustainLevel, cursor + attack + decay);
    gain.gain.setValueAtTime(peak * sustainLevel, cursor + dur - release);
    gain.gain.linearRampToValueAtTime(0, cursor + dur);
    osc.connect(gain);
    gain.connect(master);
    osc.start(cursor);
    osc.stop(cursor + dur + 0.05);
    oscillators.push(osc);

    if (opts.onNote) {
      const delay = Math.max(0, (cursor - ctx.currentTime) * 1000);
      timers.push(window.setTimeout(() => {
        if (!stopped) opts.onNote!(i, note);
      }, delay));
    }

    cursor += dur;
  });

  if (opts.onEnd) {
    const delay = Math.max(0, (cursor - ctx.currentTime) * 1000) + 30;
    timers.push(window.setTimeout(() => {
      if (!stopped) opts.onEnd!();
    }, delay));
  }

  return {
    stop() {
      if (stopped) return;
      stopped = true;
      for (const o of oscillators) {
        try {
          o.stop();
        } catch {}
      }
      for (const t of timers) clearTimeout(t);
      ctx.close().catch(() => {});
    },
    get playing() {
      return !stopped;
    },
  };
}
