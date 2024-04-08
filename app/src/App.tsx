import { useState, useEffect } from 'react'
import { Dotplots } from './Dotplots'
import { Region } from './TouchPad'
import { Config } from './Config'

function App() {
  const style = {
    width: '100vw',
    height: '100vh',
    position: 'relative',
    background: 'gray',
  } as React.CSSProperties

  const [region, setRegion] = useState<Region>({
    center: { x: 0, y: 0 },
    scale: 1,
  })
  console.log('region', region)

  return (
    <main style={style}>
      <Config region={region} />
      <Dotplots region={region} onChangeRegion={setRegion} />
    </main>
  )
}

export default App
