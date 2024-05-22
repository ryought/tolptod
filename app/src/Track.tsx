import React, { useMemo } from 'react'
import { useBoudingRect } from './useDom'

export type Feature = {
  start: number
  end: number
  seqname: string
  source: string
  type: string
  strand: string
  attributes: string
  label: string
}

const featureToLabel = (f: Feature) => {
  return `${f.seqname}:${f.start}:${f.end}:${f.strand}\n${f.source}\n${f.type}\n${f.attributes}`
}

const featureToColor = (f: Feature): [string, number] => {
  if (f.type === 'CDS') {
    return ['#ff0000', 1]
  } else if (f.type === 'exon') {
    return ['#00ff00', 2]
  } else if (f.type === 'gene') {
    return ['#0000ff', 3]
  } else {
    return ['#222222', 0]
  }
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
      [H]: 40,
    }),
    []
  )
  const width = rect[W]
  const start = center - (scale * width) / 2
  const items = useMemo(
    () =>
      features.map((feature, i) => {
        const [color, offset] = featureToColor(feature)
        const style = {
          [W]: (feature.end - feature.start) / scale,
          [H]: 10,
          [w]: (feature.start - start) / scale,
          [h]: offset * 10,
          background: color,
          opacity: 0.5,
          position: 'absolute',
          fontSize: 10,
          zIndex: 50,
          overflow: 'hidden',
        } as React.CSSProperties
        const label = feature.label || featureToLabel(feature)
        return (
          <div style={style} key={i} title={label}>
            {label}
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
