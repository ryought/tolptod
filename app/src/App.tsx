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

  const [size, setSize] = useState({ width: 0, height: 0 })
  const [region, setRegion] = useState<Region>({
    center: { x: 0, y: 0 },
    scale: 1,
  })

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
    setPlots((plots) => [...plots, plot])
  }

  return (
    <main style={style}>
      <Config region={region} onAdd={addPlot} />
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
