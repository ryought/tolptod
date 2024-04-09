import { useState, useEffect, useRef, useCallback } from 'react'
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
  // bp per px
  scale: number
}

export type Plot = {
  x: number
  y: number
  scale: number
  width: number
  height: number
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
  const [freqLow, setFreqLow] = useState(0)
  const [freqUp, setFreqUp] = useState(10)

  // touchpad related
  const [size, setSize] = useState({ width: 0, height: 0 })
  const { width, height } = size
  const [region, setRegion] = useState<Region>({
    center: { x: 0, y: 0 },
    scale: 1,
  })

  // dotplots
  const count = useRef<number>(0)
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
    const request = {
      x: targetIndex,
      y: queryIndex,
      xA,
      xB,
      yA,
      yB,
      k,
      freqLow,
      freqUp,
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
      x: xA,
      y: yA,
      scale,
      width,
      height,
      el: (
        <Dotplot
          key={count.current}
          width={width}
          height={height}
          points={points}
        />
      ),
    }
    count.current += 1
    setPlots((plots) => {
      // if (plots.length > 0) {
      //   const ret = [plots[plots.length - 1], plot]
      //   return ret
      // } else {
      //   return [plot]
      // }
      return [plot]
    })
  }

  const debounced = useDebounce(requestPlot)
  useEffect(() => {
    debounced()
  }, [region, queryIndex, targetIndex, k, freqLow, freqUp])

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
        // Id related
        targets={targets}
        querys={querys}
        targetIndex={targetIndex}
        queryIndex={queryIndex}
        onChangeTargetIndex={setTargetIndex}
        onChangeQueryIndex={setQueryIndex}
      />
      <Dotplots
        region={region}
        onChangeRegion={setRegion}
        onSizeChange={setSize}
        plots={plots}
      />
    </main>
  )
}

export default App
