import React from 'react'
import { Region } from './TouchPad'
import { Record, Plot, Job, Cache } from './App'

type Props = {
  region: Region
  onChangeRegion: (region: Region) => void
  onAdd: () => void
  onSave: () => void
  live: boolean
  onChangeLive: (live: boolean) => void
  plots: Plot[]
  onChangePlots: (plots: Plot[]) => void
  colorForward: string
  onChangeColorForward: (color: string) => void
  colorBackward: string
  onChangeColorBackward: (color: string) => void
  backgroundColor: string
  onChangeBackgroundColor: (color: string) => void
  // k-mer related
  k: number
  onChangeK: (k: number) => void
  freqLow: number
  onChangeFreqLow: (f: number) => void
  freqUp: number
  onChangeFreqUp: (f: number) => void
  localFreqLow: number
  onChangeLocalFreqLow: (f: number) => void
  localFreqUp: number
  onChangeLocalFreqUp: (f: number) => void
  dotSize: number
  onChangeDotSize: (dotSize: number) => void
  // Id related
  targetIndex: number
  queryIndex: number
  targets: Record[]
  querys: Record[]
  onChangeTargetIndex: (index: number) => void
  onChangeQueryIndex: (index: number) => void
  // cache
  cacheScale: number
  onChangeCacheScale: (scale: number) => void
  caches: Cache[]
  cacheId: string | null
  onChangeCacheId: (cacheId: string | null) => void
  onAddCache: () => void
  onRemoveCache: (cacheId: string) => void
  onUpdateCache: () => void
  // feature
  showFeature: boolean
  onChangeShowFeature: (showFeature: boolean) => void
  // job
  jobs: Job[]
  onCancelJob: (id: string) => void
}

export const Config: React.FC<Props> = ({
  region,
  onChangeRegion,
  onAdd,
  onSave,
  live,
  onChangeLive,
  plots,
  onChangePlots,
  colorForward,
  onChangeColorForward,
  colorBackward,
  onChangeColorBackward,
  backgroundColor,
  onChangeBackgroundColor,
  k,
  onChangeK,
  freqLow,
  onChangeFreqLow,
  freqUp,
  onChangeFreqUp,
  localFreqLow,
  onChangeLocalFreqLow,
  localFreqUp,
  onChangeLocalFreqUp,
  dotSize,
  onChangeDotSize,
  targets,
  querys,
  targetIndex,
  queryIndex,
  onChangeTargetIndex,
  onChangeQueryIndex,
  caches,
  cacheId,
  onChangeCacheId,
  cacheScale,
  onChangeCacheScale,
  onAddCache,
  onRemoveCache,
  onUpdateCache,
  showFeature,
  onChangeShowFeature,
  jobs,
  onCancelJob,
}) => {
  const style = {
    position: 'absolute',
    background: 'white',
    opacity: 1,
    border: 'solid',
    zIndex: 100,
    padding: 10,
    margin: 10,
    maxHeight: '80vh',
    overflow: 'scroll',
  } as React.CSSProperties
  const targetIds = targets.map(
    (record) => `${record.id} (${record.len.toLocaleString('en-US')}bp)`
  )
  const queryIds = querys.map(
    (record) => `${record.id} (${record.len.toLocaleString('en-US')}bp)`
  )
  const targetLen = targets[targetIndex]?.len || 0
  const queryLen = querys[queryIndex]?.len || 0
  const useCache = cacheId !== null

  return (
    <div style={style}>
      <details open>
        <summary></summary>
        <div>
          x(target)=
          <List
            disabled={useCache}
            items={targetIds}
            index={targetIndex}
            onChange={onChangeTargetIndex}
          />
        </div>
        <div>
          y(query)=
          <List
            disabled={useCache}
            items={queryIds}
            index={queryIndex}
            onChange={onChangeQueryIndex}
          />
        </div>
        <div>
          plot
          <button onClick={onAdd}>add</button>
          live
          <CheckBox value={live} onChange={onChangeLive} />
          {jobs.map((job) => (
            <button key={job.id} onClick={() => onCancelJob(job.id)}>
              cancel {job.id.slice(0, 5)}
            </button>
          ))}
        </div>
        <div>
          cache
          <button onClick={onAddCache}>add</button>
          <button onClick={onUpdateCache}>update</button>
          {caches.map((cache) => {
            const config = cache.config
            const summary = `X${config.x}Y${config.y} k${config.k}w${config.bin} f${config.freqLow}:${config.freqUp}`
            return (
              <div key={cache.id}>
                <CheckBox
                  disabled={!cache.done}
                  value={cacheId === cache.id}
                  onChange={(checked) => {
                    if (checked) onChangeCacheId(cache.id)
                    else onChangeCacheId(null)
                  }}
                />
                {cache.id}({cache.done ? 'done' : `${cache.progress}%`}):
                {summary}
                <button onClick={() => onRemoveCache(cache.id)}>remove</button>
              </div>
            )
          })}
        </div>
        <div>
          showFeature
          <CheckBox value={showFeature} onChange={onChangeShowFeature} />
        </div>
        <NumAndSliderInput
          label="k"
          value={k}
          onChange={onChangeK}
          min={1}
          max={200}
          disabled={useCache}
        />
        <NumAndSliderInput
          label="freqLow"
          value={freqLow}
          onChange={onChangeFreqLow}
          min={0}
          max={freqUp}
          disabled={useCache}
        />
        <NumAndSliderInput
          label="freqUp"
          value={freqUp}
          onChange={onChangeFreqUp}
          min={0}
          max={100}
          disabled={useCache}
        />
        <NumAndSliderInput
          label="localFreqLow"
          value={localFreqLow}
          onChange={onChangeLocalFreqLow}
          min={0}
          max={localFreqUp}
          disabled={useCache}
        />
        <NumAndSliderInput
          label="localFreqUp"
          value={localFreqUp}
          onChange={onChangeLocalFreqUp}
          min={0}
          max={100}
          disabled={useCache}
        />
        <NumAndSliderInput
          label="cx(bp)"
          value={region.center.x}
          onChange={(x) => {
            const newRegion = { ...region }
            newRegion.center.x = x
            onChangeRegion(newRegion)
          }}
          min={0}
          max={targetLen}
        />
        <NumAndSliderInput
          label="cy(bp)"
          value={region.center.y}
          onChange={(y) => {
            const newRegion = { ...region }
            newRegion.center.y = y
            onChangeRegion(newRegion)
          }}
          min={0}
          max={queryLen}
        />
        <div>
          scale(bp/px)
          <NumInput
            value={region.scale}
            onChange={(scale) => {
              const newRegion = { ...region }
              newRegion.scale = scale
              onChangeRegion(newRegion)
            }}
          />
        </div>
        <div>
          cache scale(bp/px)
          <NumInput value={cacheScale} onChange={onChangeCacheScale} />
        </div>
        <div>
          forward
          <input
            type="color"
            value={colorForward}
            onChange={(e) => onChangeColorForward(e.target.value)}
          />
        </div>
        <div>
          backward
          <input
            type="color"
            value={colorBackward}
            onChange={(e) => onChangeColorBackward(e.target.value)}
          />
        </div>
        <div>
          background
          <input
            type="color"
            value={backgroundColor}
            onChange={(e) => onChangeBackgroundColor(e.target.value)}
          />
        </div>
        <div>
          dotSize
          <NumInput value={dotSize} onChange={onChangeDotSize} />
        </div>
        <button onClick={onSave}>save</button>
        {plots.map((plot, i) => {
          return (
            <div key={plot.key}>
              Plot#{i}
              <CheckBox
                value={plot.active}
                onChange={(active) => {
                  const newPlots = [...plots]
                  const newPlot: Plot = {
                    ...plot,
                    active,
                  }
                  newPlots[i] = newPlot
                  onChangePlots(newPlots)
                }}
              />
              <button
                onClick={() => {
                  const newPlots = [...plots]
                  newPlots.splice(i, 1)
                  onChangePlots(newPlots)
                }}
              >
                Remove
              </button>
            </div>
          )
        })}
      </details>
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
  disabled?: boolean
}

export const Slider: React.FC<SliderProps> = ({
  value,
  onChange,
  min,
  max,
  step,
  disabled,
}) => {
  return (
    <input
      type="range"
      disabled={disabled}
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
  disabled?: boolean
}

export const List: React.FC<ListProps> = ({
  items,
  index,
  onChange,
  disabled,
}) => {
  return (
    <select
      disabled={disabled}
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

type NumInputProps = {
  value: number
  onChange: (value: number) => void
  min?: number
  max?: number
  disabled?: boolean
}

export const NumInput: React.FC<NumInputProps> = ({
  value,
  onChange,
  min,
  max,
  disabled,
}) => {
  return (
    <input
      type="number"
      disabled={disabled}
      value={value}
      min={min}
      max={max}
      style={{ width: 100 }}
      onChange={(e) => onChange(parseFloat(e.target.value))}
    />
  )
}

type NumAndSliderInputProps = {
  label: string
  value: number
  onChange: (value: number) => void
  min: number
  max: number
  disabled?: boolean
}

export const NumAndSliderInput: React.FC<NumAndSliderInputProps> = ({
  label,
  value,
  onChange,
  min,
  max,
  disabled,
}) => {
  return (
    <div>
      <div>
        {label}
        <NumInput
          value={value}
          onChange={onChange}
          min={min}
          max={max}
          disabled={disabled}
        />
      </div>
      <Slider
        value={value}
        onChange={onChange}
        min={min}
        max={max}
        disabled={disabled}
      />
    </div>
  )
}

type CheckBoxProps = {
  value: boolean
  onChange: (value: boolean) => void
  disabled?: boolean
}

const CheckBox: React.FC<CheckBoxProps> = ({ value, onChange, disabled }) => {
  return (
    <input
      type="checkbox"
      checked={value}
      disabled={disabled}
      onChange={(e) => onChange(e.target.checked)}
    />
  )
}
