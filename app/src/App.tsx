import { useState, useEffect, useRef } from 'react'
import { Dotplots } from './Dotplots'
import { Dotplot } from './Dotplot'
import { Region } from './TouchPad'
import { Config } from './Config'

export interface Record {
  id: string
  len: number
}

export interface Request {
  // regions
  x: string
  y: string
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

  useEffect(() => {
    // load query/target ids from api
    fetch('http://localhost:8080/')
      .then((res) => res.json())
      .then((json) => {
        setTargets(json['xs'] as Record[])
        setQuerys(json['ys'] as Record[])
      })
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
    const xA = Math.round(region.center.x - (width * region.scale) / 2)
    const xB = Math.round(region.center.x + (width * region.scale) / 2)
    const yA = Math.round(region.center.y - (height * region.scale) / 2)
    const yB = Math.round(region.center.y + (height * region.scale) / 2)
    const request = {
      x: targets[targetIndex].id,
      y: querys[queryIndex].id,
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
        console.log('return!', points)
        addPlot(request, points)
      })
  }
  const addPlot = (req: Request, points: [number, number][]) => {
    const { xA, xB, yA, yB, scale } = req
    const width = Math.ceil((xB - xA) / scale)
    const height = Math.ceil((yB - yA) / scale)
    const plot = {
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
      if (plots.length > 0) {
        const ret = [plots[plots.length - 1], plot]
        return ret
      } else {
        return [plot]
      }
    })
  }

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
