import React, { useState } from 'react'
import { TouchPad, Region } from './TouchPad'
import { TickBar } from './TickBar'
import { Dotplot } from './Dotplot'

type Props = {
  region: Region
  onChangeRegion: (region: Region) => void
}

type Plot = {
  x: number
  y: number
  scale: number
  width: number
  height: number
  el: JSX.Element
}

export const Dotplots: React.FC<Props> = ({ region, onChangeRegion }) => {
  const [width, setWidth] = useState(0)
  const [height, setHeight] = useState(0)

  // const [count, setCount] = useState(0)
  const style = {
    width: '100%',
    height: '100%',
  } as React.CSSProperties

  const [plots, setPlots] = useState<Plot[]>([])
  const addPlot = () => {
    const plot = {
      x: region.center.x - (width / 2) * region.scale,
      y: region.center.y - (height / 2) * region.scale,
      scale: region.scale,
      width,
      height,
      el: <Dotplot width={width} height={height} />,
    }
    setPlots((plots) => [...plots, plot])
  }

  return (
    <TouchPad
      style={style}
      region={region}
      onChange={onChangeRegion}
      onTouchEnd={() => {}}
      onSizeChange={({ width, height }) => {
        setWidth(width)
        setHeight(height)
      }}
    >
      <button onClick={() => addPlot()}>generate</button>
      <TickBar direction="x" center={region.center.x} scale={region.scale} />
      <TickBar direction="y" center={region.center.y} scale={region.scale} />
      {plots.map((plot) => {
        const style: React.CSSProperties = {
          position: 'absolute',
          left: width / 2 + (plot.x - region.center.x) / region.scale,
          top: height / 2 + (plot.y - region.center.y) / region.scale,
          width: (plot.width * plot.scale) / region.scale,
          height: (plot.height * plot.scale) / region.scale,
        }
        return <div style={style}>{plot.el}</div>
      })}
    </TouchPad>
  )
}
