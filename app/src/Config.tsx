import React, { useState } from 'react'
import { Region } from './TouchPad'
import { Record } from './App'

type Props = {
  region: Region
  onAdd: () => void
  // k-mer related
  k: number
  onChangeK: (k: number) => void
  freqLow: number
  onChangeFreqLow: (f: number) => void
  freqUp: number
  onChangeFreqUp: (f: number) => void
  // Id related
  targetIndex: number
  queryIndex: number
  targets: Record[]
  querys: Record[]
  onChangeTargetIndex: (index: number) => void
  onChangeQueryIndex: (index: number) => void
}

export const Config: React.FC<Props> = ({
  region,
  onAdd,
  k,
  onChangeK,
  freqLow,
  onChangeFreqLow,
  freqUp,
  onChangeFreqUp,
  targets,
  querys,
  targetIndex,
  queryIndex,
  onChangeTargetIndex,
  onChangeQueryIndex,
}) => {
  const style = {
    position: 'absolute',
    background: 'white',
    zIndex: 100,
    padding: 10,
    margin: 10,
  } as React.CSSProperties
  const targetIds = targets.map((record) => record.id)
  const queryIds = querys.map((record) => record.id)
  return (
    <div style={style}>
      <div>
        target
        <List
          items={targetIds}
          index={targetIndex}
          onChange={onChangeTargetIndex}
        />
        len={targets[targetIndex]?.len}
      </div>
      <div>
        query
        <List
          items={queryIds}
          index={queryIndex}
          onChange={onChangeQueryIndex}
        />
        len={querys[queryIndex]?.len}
      </div>
      <button onClick={onAdd}>add</button>
      <div>k={k}</div>
      <Slider value={k} onChange={onChangeK} min={1} max={100} />
      <div>freqLow={freqLow}</div>
      <Slider value={freqLow} onChange={onChangeFreqLow} min={1} max={freqUp} />
      <div>freqUp={freqUp}</div>
      <Slider
        value={freqUp}
        onChange={onChangeFreqUp}
        min={freqLow}
        max={100}
      />
      <div>cx={region.center.x.toFixed(0)}</div>
      <div>cy={region.center.y.toFixed(0)}</div>
      <div>scale={region.scale.toFixed(3)}</div>
    </div>
  )
}

type FileInputProps = {
  accept: string
  onLoad: (text: string) => void
}

export const FileInput: React.FC<FileInputProps> = ({ accept, onLoad }) => {
  const onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.currentTarget.files
    if (!files || files.length === 0) return
    const file = files[0]

    // read the file as text
    const reader = new FileReader()
    reader.addEventListener('load', () => {
      const text = reader.result
      if (typeof text === 'string') onLoad(text)
    })
    reader.readAsText(file, 'UTF-8')
  }
  return <input type="file" accept={accept} onChange={onChange} />
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

type ListProps = {
  items: string[]
  index: number
  onChange: (index: number) => void
}

export const List: React.FC<ListProps> = ({ items, index, onChange }) => {
  return (
    <select
      value={index.toString()}
      onChange={(e) => {
        const value = e.target.value
        onChange(parseInt(value))
      }}
    >
      {items.map((item, index) => (
        <option key={index} value={index.toString()}>
          {item}
        </option>
      ))}
    </select>
  )
}
