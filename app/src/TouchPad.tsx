import React, { useEffect, PointerEvent, WheelEvent, useRef } from 'react'
import { useBoudingRect, Rect } from './useDom'

export type Point = {
  x: number
  y: number
}

function distance(a: Point, b?: Point): number {
  if (!b) {
    return 0
  }
  return Math.sqrt((a.x - b.x) ** 2 + (a.y - b.y) ** 2)
}

function middle(a: Point, b?: Point): Point {
  if (!b) {
    return a
  }
  return {
    x: (a.x + b.x) / 2,
    y: (a.y + b.y) / 2,
  }
}

export type Region = {
  center: Point
  scale: number
}

export type Props = {
  /** movable direction (only x-axis or y-axis, or both) */
  // move: 'x' | 'y' | 'xy'
  region: Region
  onChange: (region: Region) => void
  style?: React.CSSProperties
  children?: React.ReactNode
  onTouchChange?: (touches: Touch[]) => void
  onTouchStart?: () => void
  onTouchEnd?: () => void
  onSizeChange?: (size: { width: number; height: number }) => void
}

const isMovable = (mode: 'x' | 'y' | 'xy') => {
  return {
    x: mode == 'x' || mode == 'xy',
    y: mode == 'y' || mode == 'xy',
  }
}

export type Touch = {
  point: Point
  screenPoint: Point
  pointerId: number
}

const eventToScreenPoint = (
  event: PointerEvent<HTMLDivElement> | WheelEvent<HTMLDivElement>,
  rect: Rect
): Point => {
  return {
    x: event.clientX - rect.left,
    y: event.clientY - rect.top,
  }
}

const screenPointToPointByRef = (
  screenPoint: Point,
  refScreenPoint: Point,
  refPoint: Point,
  scale: number
): Point => {
  return {
    x: refPoint.x + (screenPoint.x - refScreenPoint.x) * scale,
    y: refPoint.y + (screenPoint.y - refScreenPoint.y) * scale,
  }
}

const screenCenter = (width: number, height: number): Point => {
  return {
    x: width / 2,
    y: height / 2,
  }
}

export const screenPointToPoint = (
  screenPoint: Point,
  region: Region,
  width: number,
  height: number
): Point => {
  const { center, scale } = region
  const centerScreenPoint = screenCenter(width, height)
  return screenPointToPointByRef(screenPoint, centerScreenPoint, center, scale)
}

//
// utils related to Map<number (pointerId), Touch>
//

const getOtherTouch = (
  touches: Map<number, Touch>,
  pointerId: number
): Touch | undefined => {
  for (const [_pointerId, touch] of touches.entries()) {
    if (_pointerId !== pointerId) return touch
  }
  return
}

const toTouchList = (touches: Map<number, Touch>): Touch[] => {
  return [...touches.values()]
}

export const TouchPad: React.FC<Props> = ({
  region,
  onChange,
  onTouchChange,
  onTouchStart,
  onTouchEnd,
  onSizeChange,
  children,
  style,
}) => {
  const { ref, rect } = useBoudingRect<HTMLDivElement>()
  const { width, height } = rect
  useEffect(() => {
    if (onSizeChange) onSizeChange({ width, height })
  }, [width, height])

  const touches = useRef<Map<number, Touch>>(new Map())

  // TODO
  const onDebug = (event: PointerEvent<HTMLDivElement>) => {
    const screenPoint = eventToScreenPoint(event, rect)
    const point = screenPointToPoint(screenPoint, region, width, height)
    const { center, scale } = region
    console.log(
      event.type,
      event.pointerId,
      scale,
      // event,
      screenPoint,
      point,
      touches.current
    )
  }

  const onDown = (event: PointerEvent<HTMLDivElement>) => {
    // onDebug(event)
    if (touches.current.size < 2) {
      const screenPoint = eventToScreenPoint(event, rect)
      const point = screenPointToPoint(screenPoint, region, width, height)
      const touch = {
        point,
        screenPoint,
        pointerId: event.pointerId,
      }
      touches.current.set(event.pointerId, touch)
      if (onTouchChange) onTouchChange(toTouchList(touches.current))
      if (onTouchStart) onTouchStart()
    }
  }
  const onMove = (event: PointerEvent<HTMLDivElement>) => {
    // onDebug(event)
    const pointerId = event.pointerId
    const screenPoint = eventToScreenPoint(event, rect)

    const touch0 = touches.current.get(pointerId)
    if (!touch0) return // ignore pointer of which pointerdown event was missed
    const touch1 = getOtherTouch(touches.current, pointerId)

    // update scale so that dist [m] corresponds to screenDist [px]
    const dist = distance(touch0.point, touch1?.point)
    const screenDist = distance(screenPoint, touch1?.screenPoint)
    const scale = dist / screenDist || region.scale

    // update region.center so that p corresponds to P
    const midPoint = middle(touch0.point, touch1?.point)
    const midScreenPoint = middle(screenPoint, touch1?.screenPoint)
    const center = screenPointToPointByRef(
      screenCenter(width, height),
      midScreenPoint,
      midPoint,
      scale
    )
    onChange({ center, scale })

    // update screenPoint
    touch0.screenPoint = screenPoint
    touches.current.set(pointerId, touch0)

    if (onTouchChange) onTouchChange(toTouchList(touches.current))
  }
  const onUp = (event: PointerEvent<HTMLDivElement>) => {
    // onDebug(event)

    touches.current.delete(event.pointerId)
    touches.current.forEach((touch) => {
      // update other touch points
      touch.point = screenPointToPoint(touch.screenPoint, region, width, height)
    })

    if (onTouchChange) onTouchChange(toTouchList(touches.current))
    if (onTouchEnd) onTouchEnd()
  }
  const onWheel = (event: WheelEvent<HTMLDivElement>) => {
    // onDebug(event)
    const screenPoint = eventToScreenPoint(event, rect)
    const point = screenPointToPoint(screenPoint, region, width, height)
    const scale = region.scale * Math.exp(event.deltaY * 0.01)

    const center = screenPointToPointByRef(
      screenCenter(width, height),
      screenPoint,
      point,
      scale
    )
    onChange({ center, scale })
  }

  return (
    <div
      ref={ref}
      style={style}
      onPointerDown={onDown}
      onPointerMove={onMove}
      onPointerUp={onUp}
      onPointerCancel={onUp}
      // TODO
      // onPointerOut={onUp}
      onPointerLeave={onUp}
      onWheel={onWheel}
    >
      {children}
    </div>
  )
}
