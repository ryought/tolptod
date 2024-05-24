import React from 'react'
import { Region } from './TouchPad'
import { Record, Plot, Job } from './App'

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
  useCache: boolean
  onChangeUseCache: (useCache: boolean) => void
  onUpdateCache: () => void
  cacheJob: Job | undefined
  onCancelCacheJob: () => void
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
  useCache,
  cacheScale,
  onChangeCacheScale,
  onChangeUseCache,
  onUpdateCache,
  showFeature,
  onChangeShowFeature,
  jobs,
  onCancelJob,
  cacheJob,
  onCancelCacheJob,
}) => {
  const style = {
    position: 'absolute',
    background: 'white',
    opacity: 1,
    border: 'solid',
    zIndex: 100,
    padding: 10,
    margin: 10,
  } as React.CSSProperties
  const targetIds = targets.map(
    (record) => `${record.id} (${record.len.toLocaleString('en-US')}bp)`
  )
  const queryIds = querys.map(
    (record) => `${record.id} (${record.len.toLocaleString('en-US')}bp)`
  )

  return (
    <div style={style}>
      <details open>
        <summary></summary>
        <div>
          x(target)=
          <List
            items={targetIds}
            index={targetIndex}
            onChange={onChangeTargetIndex}
          />
        </div>
        <div>
          y(query)=
          <List
            items={queryIds}
            index={queryIndex}
            onChange={onChangeQueryIndex}
          />
        </div>
        <div>
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
          {cacheJob ? (
            <button onClick={onCancelCacheJob}>cancel</button>
          ) : (
            <button onClick={onUpdateCache}>update cache</button>
          )}
          useCache
          <CheckBox value={useCache} onChange={onChangeUseCache} />
          <button
            onClick={() => {
              fetch('http://localhost:8080/cache/')
                .then((res) => res.json())
                .then((json) => console.log('get /cache', json))
                .catch((err) => console.error(err))
            }}
          >
            List
          </button>
        </div>
        <div>
          showFeature
          <CheckBox value={showFeature} onChange={onChangeShowFeature} />
        </div>
        <div>
          k
          <NumInput value={k} onChange={onChangeK} />
        </div>
        <Slider value={k} onChange={onChangeK} min={1} max={100} />
        <NumAndSliderInput
          label="freqLow"
          value={freqLow}
          onChange={onChangeFreqLow}
          min={0}
          max={freqUp}
        />
        <NumAndSliderInput
          label="freqUp"
          value={freqUp}
          onChange={onChangeFreqUp}
          min={0}
          max={100}
        />
        <NumAndSliderInput
          label="localFreqLow"
          value={localFreqLow}
          onChange={onChangeLocalFreqLow}
          min={0}
          max={localFreqUp}
        />
        <NumAndSliderInput
          label="localFreqUp"
          value={localFreqUp}
          onChange={onChangeLocalFreqUp}
          min={0}
          max={100}
        />
        <div>
          cx(bp)
          <NumInput
            value={region.center.x}
            onChange={(x) => {
              const newRegion = { ...region }
              newRegion.center.x = x
              onChangeRegion(newRegion)
            }}
          />
        </div>
        <div>
          cy(bp)
          <NumInput
            value={region.center.y}
            onChange={(y) => {
              const newRegion = { ...region }
              newRegion.center.y = y
              onChangeRegion(newRegion)
            }}
          />
        </div>
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

type NumInputProps = {
  value: number
  onChange: (value: number) => void
  min?: number
  max?: number
}

export const NumInput: React.FC<NumInputProps> = ({
  value,
  onChange,
  min,
  max,
}) => {
  return (
    <input
      type="number"
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
}

export const NumAndSliderInput: React.FC<NumAndSliderInputProps> = ({
  label,
  value,
  onChange,
  min,
  max,
}) => {
  return (
    <div>
      <div>
        {label}
        <NumInput value={value} onChange={onChange} min={min} max={max} />
      </div>
      <Slider value={value} onChange={onChange} min={min} max={max} />
    </div>
  )
}

type CheckBoxProps = {
  value: boolean
  onChange: (value: boolean) => void
}

const CheckBox: React.FC<CheckBoxProps> = ({ value, onChange }) => {
  return (
    <input
      type="checkbox"
      checked={value}
      onChange={(e) => onChange(e.target.checked)}
    />
  )
}
