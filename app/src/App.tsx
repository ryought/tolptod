import { useState, useEffect } from 'react'
import { Dotplots } from './Dotplots'
import { Dotplot, Plot } from './Dotplot'
import { Region } from './TouchPad'
import { Config } from './Config'

function App() {
  const style = {
    width: '100vw',
    height: '100vh',
    position: 'relative',
    background: 'gray',
  } as React.CSSProperties

  // sequence names
  const [queryIds, setQueryIds] = useState<string[]>([])
  const [targetIds, setTargetIds] = useState<string[]>(['hogehoge'])
  const [query, setQuery] = useState<number>(0)
  const [target, setTarget] = useState<number>(0)

  useEffect(() => {
    // load ids from api
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
        onAdd={addPlot}
        // k-mer rleated
        k={k}
        onChangeK={setK}
        freqLow={freqLow}
        onChangeFreqLow={setFreqLow}
        freqUp={freqUp}
        onChangeFreqUp={setFreqUp}
        // Id related
        targetIds={targetIds}
        queryIds={queryIds}
        target={target}
        query={query}
        onChangeTarget={setTarget}
        onChangeQuery={setQuery}
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
