import { useState, useEffect, useRef } from 'react'
import { Dotplots } from './Dotplots'
import { Dotplot, Points } from './Dotplot'
import { Region } from './TouchPad'
import type { Feature } from './Track'
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
  // cache
  useCache: boolean
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

export const floorMultiple = (x: number, m: number) => {
  return Math.floor(x / m) * m
}

function App() {
  // sequence names
  const [querys, setQuerys] = useState<Record[]>([])
  const [targets, setTargets] = useState<Record[]>([])
  const [queryIndex, setQueryIndex] = useState<number>(0)
  const [targetIndex, setTargetIndex] = useState<number>(0)
  const queryLen = querys[queryIndex]?.len || 0
  const targetLen = targets[targetIndex]?.len || 0
  const isDev = import.meta.env.MODE === 'development'

  useEffect(() => {
    // load query/target ids from api
    fetch(isDev ? 'http://localhost:8080/info/' : 'info/')
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
  const [freqUp, setFreqUp] = useState(1)

  // touchpad related
  const [size, setSize] = useState({ width: 0, height: 0 })
  const { width, height } = size
  const [region, setRegion] = useState<Region>({
    center: { x: 0, y: 0 },
    scale: 1,
  })

  // feature
  const [features, setFeatures] = useState<{ x: Feature[]; y: Feature[] }>({
    x: [],
    y: [],
  })

  // dotplots
  const [colorForward, setColorForward] = useState<string>('#FF0000')
  const [colorBackward, setColorBackward] = useState<string>('#0000FF')
  const [backgroundColor, setBackgroundColor] = useState<string>('#808080')
  const count = useRef<number>(0)
  const [live, setLive] = useState<boolean>(true)
  const [useCache, setUseCache] = useState<boolean>(false)
  const [showFeature, setShowFeature] = useState<boolean>(true)
  const [currentPlot, setCurrentPlot] = useState<Plot | null>(null)
  const [plots, setPlots] = useState<Plot[]>([])
  const scale = Math.ceil(region.scale)
  const [cacheScale, setCacheScale] = useState<number>(1)
  useEffect(() => {
    setCacheScale(
      Math.max(
        Math.pow(2, Math.floor(Math.log2(queryLen || 1) / 2)),
        Math.pow(2, Math.floor(Math.log2(targetLen || 1) / 2))
      )
    )
  }, [querys, targets, queryIndex, targetIndex])
  const requestPlot = () => {
    const data = new FormData()
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
    const s = useCache
      ? Math.max(1, Math.floor(scale / cacheScale)) * cacheScale
      : scale
    const request: Request = {
      x: targetIndex,
      y: queryIndex,
      xA: useCache ? floorMultiple(xA, s) : xA,
      xB: useCache ? floorMultiple(xB, s) : xB,
      yA: useCache ? floorMultiple(yA, s) : yA,
      yB: useCache ? floorMultiple(yB, s) : yB,
      k,
      scale: s,
      freqLow,
      freqUp,
      useCache,
    }
    data.append('json', JSON.stringify(request))
    console.log('send', request)
    fetch(isDev ? 'http://localhost:8080/generate/' : 'generate/', {
      method: 'POST',
      body: data,
    })
      .then((res) => res.json())
      .then((json) => {
        const points = {
          forward: json.forward as [number, number][],
          backward: json.backward as [number, number][],
        }
        // console.log('points', points)
        addPlot(request, points)
      })
      .catch(() => alert('cannot /generate'))
    if (showFeature) {
      fetch(isDev ? 'http://localhost:8080/features/' : 'features/', {
        method: 'POST',
        body: data,
      })
        .then((res) => res.json())
        .then((json) => {
          setFeatures({
            x: json.x ? json.x : [],
            y: json.y ? json.y : [],
          })
          console.log('features', json)
        })
    }
  }
  const addPlot = (req: Request, points: Points) => {
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
          colorForward={colorForward}
          colorBackward={colorBackward}
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
  }, [
    live,
    region,
    queryIndex,
    targetIndex,
    k,
    freqLow,
    freqUp,
    colorForward,
    colorBackward,
  ])

  const onUpdateCache = () => {
    const data = new FormData()
    const request: Request = {
      x: targetIndex,
      y: queryIndex,
      xA: 0,
      xB: targetLen,
      yA: 0,
      yB: queryLen,
      k,
      freqLow,
      freqUp,
      scale: cacheScale,
      useCache,
    }
    data.append('json', JSON.stringify(request))
    fetch(isDev ? 'http://localhost:8080/cache/' : 'cache/', {
      method: 'POST',
      body: data,
    })
  }

  const style = {
    width: '100vw',
    height: '100vh',
    position: 'relative',
    background: backgroundColor,
  } as React.CSSProperties

  return (
    <main style={style}>
      <Config
        region={region}
        onChangeRegion={setRegion}
        onAdd={requestPlot}
        // k-mer rleated
        k={k}
        onChangeK={setK}
        freqLow={freqLow}
        onChangeFreqLow={setFreqLow}
        freqUp={freqUp}
        onChangeFreqUp={setFreqUp}
        // color
        colorForward={colorForward}
        onChangeColorForward={setColorForward}
        colorBackward={colorBackward}
        onChangeColorBackward={setColorBackward}
        backgroundColor={backgroundColor}
        onChangeBackgroundColor={setBackgroundColor}
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
        onUpdateCache={onUpdateCache}
        useCache={useCache}
        onChangeUseCache={setUseCache}
        cacheScale={cacheScale}
        onChangeCacheScale={setCacheScale}
        showFeature={showFeature}
        onChangeShowFeature={setShowFeature}
      />
      <Dotplots
        region={region}
        onChangeRegion={setRegion}
        onSizeChange={setSize}
        plots={currentPlot ? [...plots, currentPlot] : plots}
        features={features}
      />
    </main>
  )
}

export default App
