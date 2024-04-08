import React from 'react'
import { useBoudingRect } from './useDom'

type Props = {
  direction: 'x' | 'y'
  // bp
  center: number
  // bp per px
  scale: number
}

export const TickBar: React.FC<Props> = ({ direction, center, scale }) => {
  const { ref, rect } = useBoudingRect<HTMLDivElement>()

  const W = direction == 'x' ? 'width' : 'height'
  const H = direction == 'x' ? 'height' : 'width'
  const w = direction == 'x' ? 'left' : 'top'
  const h = direction == 'x' ? 'top' : 'left'
  const style: React.CSSProperties = {
    [W]: '100%',
  }
  const width = rect[W]

  const total = width * scale
  const interval = Math.pow(10, Math.floor(Math.log10(total) - 0.3))
  const start = center - (scale * width) / 2
  const firstTick = Math.floor(start / interval)
  const nTicks = Math.ceil(total / interval)
  const tickPositions = Array.from(
    { length: nTicks },
    (_, i) => (firstTick + i) * interval
  )

  const ticks = tickPositions.map((t) => {
    const style = {
      [W]: 1,
      [H]: 20,
      [w]: (t - start) / scale,
      [h]: 0,
      background: 'black',
      position: 'absolute',
      fontSize: 10,
      zIndex: 50,
    } as React.CSSProperties
    return (
      <div style={style} key={t}>
        {t.toLocaleString()}
      </div>
    )
  })

  return (
    <div ref={ref} style={style}>
      {ticks}
    </div>
  )
}
