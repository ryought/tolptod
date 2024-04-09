import { useState, useEffect } from 'react'
import { Dotplots } from './Dotplots'
import { Dotplot, Plot } from './Dotplot'
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
  const [region, setRegion] = useState<Region>({
    center: { x: 0, y: 0 },
    scale: 1,
  })

  // dotplots
  const [plots, setPlots] = useState<Plot[]>([])
  const requestPlot = () => {
    const request = {
      x: targets[targetIndex].id,
      y: querys[queryIndex].id,
      k,
      freqLow,
      freqUp,
      // scale,
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
        console.log('return!', json)
      })
  }
  const addPlot = () => {
    const plot = {
      x: region.center.x - (size.width / 2) * region.scale,
      y: region.center.y - (size.height / 2) * region.scale,
      scale: region.scale,
      width: size.width,
      height: size.height,
      el: <Dotplot width={size.width} height={size.height} />,
    }
    setPlots((plots) => {
      const ret = [...plots.slice(-2), plot]
      console.log('ret', ret)
      return ret
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
