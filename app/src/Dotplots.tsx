import React, { useState } from 'react'
import { TouchPad, Region } from './TouchPad'
import { TickBar } from './TickBar'
import { Track, Feature } from './Track'
import { Plot } from './App'

type Props = {
  region: Region
  onChangeRegion: (region: Region) => void
  onSizeChange: (size: { width: number; height: number }) => void
  onTouchEnd?: () => void
  plots: Plot[]
  features: Feature[]
}

export const Dotplots: React.FC<Props> = ({
  region,
  onChangeRegion,
  onSizeChange,
  onTouchEnd = () => {},
  plots,
  features,
}) => {
  const [size, setSize] = useState({ width: 0, height: 0 })
  const style = {
    width: '100%',
    height: '100%',
  } as React.CSSProperties
  return (
    <TouchPad
      style={style}
      region={region}
      onChange={onChangeRegion}
      onTouchEnd={onTouchEnd}
      onSizeChange={(size) => {
        setSize(size)
        onSizeChange(size)
      }}
    >
      <Track
        direction="x"
        center={region.center.x}
        scale={region.scale}
        features={features}
      />
      <Track
        direction="y"
        center={region.center.y}
        scale={region.scale}
        features={features}
      />
      <TickBar direction="x" center={region.center.x} scale={region.scale} />
      <TickBar direction="y" center={region.center.y} scale={region.scale} />
      {plots
        .filter((plot) => plot.active)
        .map((plot, i) => {
          const style: React.CSSProperties = {
            position: 'absolute',
            left: size.width / 2 + (plot.x - region.center.x) / region.scale,
            top: size.height / 2 + (plot.y - region.center.y) / region.scale,
            width: (plot.width * plot.scale) / region.scale,
            height: (plot.height * plot.scale) / region.scale,
          }
          return (
            <div key={i} style={style}>
              {plot.el}
            </div>
          )
        })}
    </TouchPad>
  )
}
