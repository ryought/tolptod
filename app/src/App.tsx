import { useState, useEffect, useRef } from 'react'
import { Dotplots } from './Dotplots'
import { Dotplot } from './Dotplot'
import { Region } from './TouchPad'
import { Config } from './Config'
import { useDebounce } from './debounce'
import { clamp } from './utils'

export interface Record {
  id: string
  len: number
}

export interface Request {
  // regions
  x: number
  y: number
  xA: number
  xB: number
  yA: number
  yB: number
  // parameters
  k: number
  freqLow: number
  freqUp: number
  revcomp: boolean
  // bp per px
  scale: number
}

export type Plot = {
  key: number
  x: number
  y: number
  scale: number
  width: number
  height: number
  active: boolean
  el: JSX.Element
}

function App() {
  const style = {
    width: '100vw',
    height: '100vh',
    position: 'relative',
    background: 'gray',
  } as React.CSSProperties

  // sequence names
  const [querys, setQuerys] = useState<Record[]>([])
  const [targets, setTargets] = useState<Record[]>([])
  const [queryIndex, setQueryIndex] = useState<number>(0)
  const [targetIndex, setTargetIndex] = useState<number>(0)
  const queryLen = querys[queryIndex]?.len || 0
  const targetLen = targets[targetIndex]?.len || 0

  useEffect(() => {
    // load query/target ids from api
    fetch('http://localhost:8080/')
      .then((res) => res.json())
      .then((json) => {
        setTargets(json['xs'] as Record[])
        setQuerys(json['ys'] as Record[])
      })
      .catch(() => alert('cannot get info'))
  }, [])

  // k-mer related
  const [k, setK] = useState(16)
  const [freqLow, setFreqLow] = useState(1)
  const [freqUp, setFreqUp] = useState(50)
  const [revcomp, setRevcomp] = useState<boolean>(false)

  // touchpad related
  const [size, setSize] = useState({ width: 0, height: 0 })
  const { width, height } = size
  const [region, setRegion] = useState<Region>({
    center: { x: 0, y: 0 },
    scale: 1,
  })

  // dotplots
  const [color, setColor] = useState<string>('#FF0000')
  const count = useRef<number>(0)
  const [live, setLive] = useState<boolean>(true)
  const [currentPlot, setCurrentPlot] = useState<Plot | null>(null)
  const [plots, setPlots] = useState<Plot[]>([])
  const requestPlot = () => {
    const scale = Math.ceil(region.scale)
    const xA = clamp(
      Math.round(region.center.x - (width * region.scale) / 2),
      0,
      targetLen
    )
    const xB = clamp(
      Math.round(region.center.x + (width * region.scale) / 2),
      0,
      targetLen
    )
    const yA = clamp(
      Math.round(region.center.y - (height * region.scale) / 2),
      0,
      queryLen
    )
    const yB = clamp(
      Math.round(region.center.y + (height * region.scale) / 2),
      0,
      queryLen
    )
    const request: Request = {
      x: targetIndex,
      y: queryIndex,
      xA,
      xB,
      yA,
      yB,
      k,
      freqLow,
      freqUp,
      revcomp,
      scale,
    }
    const data = new FormData()
    data.append('json', JSON.stringify(request))
    console.log('send', request)
    fetch('http://localhost:8080/generate/', {
      method: 'POST',
      body: data,
    })
      .then((res) => res.json())
      .then((json) => {
        const points = json.points as [number, number][]
        addPlot(request, points)
      })
      .catch(() => alert('cannot /generate'))
  }
  const addPlot = (req: Request, points: [number, number][]) => {
    const { xA, xB, yA, yB, scale } = req
    const width = Math.ceil((xB - xA) / scale)
    const height = Math.ceil((yB - yA) / scale)
    const plot: Plot = {
      key: count.current,
      x: xA,
      y: yA,
      scale,
      width,
      height,
      active: true, // new plot is always active
      el: (
        <Dotplot
          key={count.current}
          width={width}
          height={height}
          points={points}
          color={color}
        />
      ),
    }
    count.current += 1
    setCurrentPlot(plot)
  }
  const savePlot = () => {
    if (currentPlot) {
      setPlots((plot) => [...plot, currentPlot])
      setCurrentPlot(null)
    }
  }

  const debounced = useDebounce(requestPlot)
  useEffect(() => {
    if (live) debounced()
  }, [live, region, queryIndex, targetIndex, k, freqLow, freqUp, color])

  return (
    <main style={style}>
      <Config
        region={region}
        onAdd={requestPlot}
        // k-mer rleated
        k={k}
        onChangeK={setK}
        freqLow={freqLow}
        onChangeFreqLow={setFreqLow}
        freqUp={freqUp}
        onChangeFreqUp={setFreqUp}
        // color
        color={color}
        onChangeColor={setColor}
        // revcomp
        revcomp={revcomp}
        onChangeRevcomp={setRevcomp}
        // Id related
        targets={targets}
        querys={querys}
        targetIndex={targetIndex}
        queryIndex={queryIndex}
        onChangeTargetIndex={setTargetIndex}
        onChangeQueryIndex={setQueryIndex}
        // save related
        onSave={savePlot}
        live={live}
        onChangeLive={setLive}
        plots={plots}
        onChangePlots={setPlots}
      />
      <Dotplots
        region={region}
        onChangeRegion={setRegion}
        onSizeChange={setSize}
        plots={currentPlot ? [...plots, currentPlot] : plots}
      />
    </main>
  )
}

export default App
