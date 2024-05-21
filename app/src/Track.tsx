import React, { useMemo } from 'react'
import { useBoudingRect } from './useDom'

export type Feature = {
  start: number
  end: number
  label: string
}

type Props = {
  // labeled regions in bp
  features: Feature[]
  //
  direction: 'x' | 'y'
  // bp
  center: number
  // bp per px
  scale: number
}

export const Track: React.FC<Props> = ({
  features,
  direction,
  center,
  scale,
}) => {
  const { ref, rect } = useBoudingRect<HTMLDivElement>()
  const W = direction == 'x' ? 'width' : 'height'
  const H = direction == 'x' ? 'height' : 'width'
  const w = direction == 'x' ? 'left' : 'top'
  const h = direction == 'x' ? 'top' : 'left'
  const style: React.CSSProperties = useMemo(
    () => ({
      [W]: '100%',
      [h]: 20,
      background: 'white',
      position: 'absolute',
      [H]: 20,
    }),
    []
  )
  const width = rect[W]
  const start = center - (scale * width) / 2
  const items = useMemo(
    () =>
      features.map((feature, i) => {
        const style = {
          [W]: (feature.end - feature.start) / scale,
          [H]: 20,
          [w]: (feature.start - start) / scale,
          [h]: 0,
          background: 'blue',
          position: 'absolute',
          fontSize: 10,
          zIndex: 50,
          overflow: 'hidden',
        } as React.CSSProperties
        return (
          <div style={style} key={i} title={feature.label}>
            {feature.label}
          </div>
        )
      }),
    [features, scale, start]
  )
  return (
    <div ref={ref} style={style}>
      {items}
    </div>
  )
}
