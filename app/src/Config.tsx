import React, { useState } from 'react'
import { Region } from './TouchPad'

type Props = {
  region: Region
  k: number
  f: number
}

export const Config: React.FC<Props> = ({ region }) => {
  const style = {
    position: 'absolute',
    background: 'white',
    zIndex: 100,
    padding: 10,
    margin: 10,
  } as React.CSSProperties
  const [k, setK] = useState(16)
  const [freqLow, setFreqLow] = useState(0)
  const [freqUp, setFreqUp] = useState(10)
  return (
    <div style={style}>
      <div>k={k}</div>
      <Slider value={k} onChange={setK} min={1} max={100} />
      <div>freqLow={freqLow}</div>
      <Slider value={freqLow} onChange={setFreqLow} min={1} max={freqUp} />
      <div>freqUp={freqUp}</div>
      <Slider value={freqUp} onChange={setFreqUp} min={freqLow} max={100} />
      <div>cx={region.center.x.toFixed(0)}</div>
      <div>cy={region.center.y.toFixed(0)}</div>
    </div>
  )
}

type SliderProps = {
  value: number
  min: number
  max: number
  step?: number
  onChange: (value: number) => void
}

export const Slider: React.FC<SliderProps> = ({
  value,
  onChange,
  min,
  max,
  step,
}) => {
  return (
    <input
      type="range"
      min={min}
      max={max}
      step={step ?? 1}
      value={value}
      onChange={(e) => onChange(parseFloat(e.target.value))}
    />
  )
}
