package gokit

import (
	"fmt"
	"math"
)

type VectorF struct {
	X float32
	Y float32
}

type point struct {
	x int
	y int
}
type pathPoint struct {
	point

	weight       int
	fillWeight   int
	distTraveled int

	parent *pathPoint
	next   *pathPoint
}

type pointList struct {
	head *pathPoint
}

func newPointList() *pointList {
	return &pointList{}
}

func (lst *pointList) push(point *pathPoint) {
	if lst.head == nil {
		lst.head = point
		return
	}

	node := lst.head
	var preNode *pathPoint
	for nil != node {
		if node.weight >= point.weight {
			if nil == preNode {
				lst.head = point
				point.next = node
			} else {
				preNode.next = point
				point.next = node
			}
			return
		}

		preNode = node
		node = node.next
	}

	preNode.next = point
}

func (lst *pointList) pop() *pathPoint {
	if nil == lst.head {
		return nil
	}

	head := lst.head
	lst.head = lst.head.next

	return head
}

type AMap struct {
	maxX        int
	maxY        int
	filledTiles map[point]int
}

//maxX maxY  x y最大值
func NewAMap(maxX, maxY int) *AMap {
	amap := &AMap{maxX: maxX, maxY: maxY}
	amap.filledTiles = make(map[point]int)

	return amap
}

//weight -1 不可移动的点
func (am *AMap) Filled(x, y, weight int) {
	am.filledTiles[point{x: x, y: y}] = weight
}

func (am *AMap) canMove(x, y int) bool {
	return am.filledTiles[point{x: x, y: y}] >= 0
}

type AStar struct {
	searchRange int
	bSmooth     bool
}

func NewAStar() *AStar {
	return &AStar{searchRange: 0, bSmooth: true}
}

func (a *AStar) getRange(pos, maxPos, iRange int) (int, int) {
	min := pos - iRange
	if min < 0 {
		min = 0
	}

	max := pos + iRange
	if max >= maxPos {
		max = maxPos - 1
	}

	return min, max
}

func (a *AStar) FindLatelyPoint(x1, y1 int, aMap *AMap) (newX, newY int, ok bool) {
	var tmpPoint []*point
	var minDis float64
	bFirst := true
	var latelyPoint *point

	for i := 1; i <= a.searchRange; i++ {
		//左
		x := x1 - i
		if x >= 0 {
			yMin, yMax := a.getRange(y1, aMap.maxY, i)
			for y := yMin; y <= yMax; y++ {
				if aMap.canMove(x, y) {
					tmpPoint = append(tmpPoint, &point{x: x, y: y})
				}
			}
		}
		//右
		x = x1 + i
		if x < aMap.maxX {
			yMin, yMax := a.getRange(y1, aMap.maxY, i)
			for y := yMin; y <= yMax; y++ {
				if aMap.canMove(x, y) {
					tmpPoint = append(tmpPoint, &point{x: x, y: y})
				}
			}
		}
		//上
		y := y1 + i
		if y < aMap.maxY {
			xMin, xMax := a.getRange(x1, aMap.maxX, i)
			for x = xMin; x <= xMax; x++ {
				if aMap.canMove(x, y) {
					tmpPoint = append(tmpPoint, &point{x: x, y: y})
				}
			}
		}
		//下
		y = y1 - i
		if y >= 0 {
			xMin, xMax := a.getRange(x1, aMap.maxX, i)
			for x = xMin; x <= xMax; x++ {
				if aMap.canMove(x, y) {
					tmpPoint = append(tmpPoint, &point{x: x, y: y})
				}
			}
		}

		//找出最近的点
		if 1 == len(tmpPoint) {
			return tmpPoint[0].x, tmpPoint[0].y, true
		}

		if len(tmpPoint) > 0 {
			for _, p := range tmpPoint {
				dis := math.Abs(float64(x1-p.x)) + math.Abs(float64(y1-p.y))
				if bFirst {
					minDis = dis
					latelyPoint = p
					bFirst = false
				}
				if dis < minDis {
					minDis = dis
					latelyPoint = p
				}
			}

			return latelyPoint.x, latelyPoint.y, true
		}
	}

	return 0, 0, false
}

func (a *AStar) correctPoint(x1, y1, x2, y2 float32, aMap *AMap) (source, target point, bAddStart, bAddEnd bool, err error) {
	bAddStart = false
	bAddEnd = false
	//四舍入五找最近的格子
	ix1 := a.clampInt(int(a.round(x1)), 0, aMap.maxX-1)
	iy1 := a.clampInt(int(a.round(y1)), 0, aMap.maxY-1)
	ix2 := a.clampInt(int(a.round(x2)), 0, aMap.maxX-1)
	iy2 := a.clampInt(int(a.round(y2)), 0, aMap.maxY-1)

	if !aMap.canMove(ix1, iy1) {
		newX, newY, ok := a.FindLatelyPoint(ix1, iy1, aMap)
		if !ok {
			err = fmt.Errorf("start point unreachable.")
			return
		}

		bAddStart = true
		source.x = newX
		source.y = newY
	} else {
		source.x = ix1
		source.y = iy1
	}

	if !aMap.canMove(ix2, iy2) {
		newX, newY, ok := a.FindLatelyPoint(ix2, iy2, aMap)
		if !ok {
			err = fmt.Errorf("end point unreachable.")
			return
		}

		bAddEnd = true
		target.x = newX
		target.y = newY
	} else {
		target.x = ix2
		target.y = iy2
	}

	return
}

func (a *AStar) setWeight(p *pathPoint, fill_weight int, end point) bool {
	if -1 == fill_weight {
		return false
	}

	p.weight = p.fillWeight + p.distTraveled + int(math.Abs(float64(p.point.x-end.x))+math.Abs(float64(p.point.y-end.y)))
	return true
}

func (a *AStar) getSurrounding(p point, aMap *AMap) []point {
	var surrounding []point
	x, y := p.x, p.y

	if x > 0 {
		surrounding = append(surrounding, point{x - 1, y})
	}
	if x < aMap.maxX-1 {
		surrounding = append(surrounding, point{x + 1, y})
	}

	if y > 0 {
		surrounding = append(surrounding, point{x, y - 1})
	}
	if y < aMap.maxY-1 {
		surrounding = append(surrounding, point{x, y + 1})
	}

	return surrounding
}

func (a *AStar) find(source, target point, aMap *AMap) []*point {
	openList := make(map[point]*pathPoint)
	closeList := make(map[point]*pathPoint)
	listOpen := newPointList()

	target_weight := aMap.filledTiles[target]
	target_point := &pathPoint{
		point:        target,
		parent:       nil,
		distTraveled: 0,
		fillWeight:   target_weight,
	}

	if a.setWeight(target_point, target_weight, source) {
		openList[target] = target_point
		listOpen.push(target_point)
	}

	var current *pathPoint
	for {
		current = listOpen.pop()
		if current == nil || current.point == source {
			break
		}

		delete(openList, current.point)
		closeList[current.point] = current

		surrounding := a.getSurrounding(current.point, aMap)
		for _, p := range surrounding {
			_, ok := closeList[p]
			if ok {
				continue
			}
			if !aMap.canMove(p.x, p.y) {
				continue
			}

			fill_weight := aMap.filledTiles[p]
			path_point := &pathPoint{
				point:        p,
				parent:       current,
				fillWeight:   current.fillWeight + fill_weight,
				distTraveled: current.distTraveled + 1,
			}
			a.setWeight(path_point, fill_weight, source)

			existing_point, ok := openList[p]
			if !ok {
				openList[p] = path_point
				listOpen.push(path_point)
			} else {
				if path_point.weight < existing_point.weight {
					existing_point.parent = path_point.parent
				}
			}
		}
	}

	if nil == current {
		return nil
	}

	var path []*point
	p := current
	for p != nil {
		path = append(path, &p.point)
		p = p.parent
	}

	return path
}

func (a *AStar) checkInLine(point1, point2, point3 *point) bool {
	if point1.x == point2.x && point1.x == point3.x {
		return true
	}
	if point1.y == point2.y && point1.y == point3.y {
		return true
	}

	return false
}

//计算起始点到终止点的路径上是否有移动阻挡，无宽度
func (a *AStar) DetectMoveCollisionBetween(startX, startY, endX, endY float32, aMap *AMap) bool {
	x0 := startX
	y0 := startY
	x1 := endX
	y1 := endY

	steep := math.Abs(float64(y1-y0)) > math.Abs(float64(x1-x0))

	if steep {
		x0 = startY
		y0 = startX
		x1 = endY
		y1 = endX
	}

	if x0 > x1 {
		x := x0
		y := y0
		x0 = x1
		x1 = x
		y0 = y1
		y1 = y
	}

	ratio := math.Abs(float64((y1 - y0) / (x1 - x0)))

	var mirror int
	if y1 > y0 {
		mirror = 1
	} else {
		mirror = -1
	}

	skip := false
	for col := int(math.Floor(float64(x0))); col < int(math.Ceil(float64(x1))); col++ {
		curY := y0 + float32(mirror)*float32(ratio)*(float32(col)-x0)
		//第一格不进行延边计算
		skip = false

		if col == int(math.Floor(float64(x0))) {
			skip = int(curY) != int(y0)
		}

		if !skip {
			if !steep {
				if !aMap.canMove(col, int(math.Max(0, math.Floor(float64(curY))))) {
					return true
				}
			} else {
				if !aMap.canMove(int(math.Max(0, math.Floor(float64(curY)))), col) {
					return true
				}
			}
		}

		var tmp float32
		if mirror > 0 {
			tmp = float32(math.Ceil(float64(curY))) - curY
		} else {
			tmp = curY - float32(math.Floor(float64(curY)))
		}

		//根据斜率计算是否有跨格
		if tmp < float32(ratio) {
			crossY := int(math.Floor(float64(curY))) + mirror
			//判断是否超出范围
			if crossY > int(math.Max(float64(y0), float64(y1))) ||
				crossY < int(math.Min(float64(y0), float64(y1))) {
				return false
			}

			//跨线格子
			if !steep {
				if !aMap.canMove(col, crossY) {
					return true
				}
			} else {
				if !aMap.canMove(crossY, col) {
					return true
				}
			}
		}
	}

	return false
}

func (a *AStar) removePoint(path []*point, aMap *AMap) []*point {
	if len(path) < 3 {
		return path
	}

	var tmpLine []*point
	tmpLine = append(tmpLine, path[0])
	//移除直线上的点
	for i := 1; i < len(path)-1; i++ {
		if !a.checkInLine(tmpLine[len(tmpLine)-1], path[i], path[i+1]) {
			tmpLine = append(tmpLine, path[i])
		}
	}
	tmpLine = append(tmpLine, path[len(path)-1])
	if len(tmpLine) < 3 {
		return tmpLine
	}

	//拐点移除
	var tmpPath []*point
	tmpPath = append(tmpPath, tmpLine[0])
	for i := 2; i < len(tmpLine); i++ {
		if a.DetectMoveCollisionBetween(float32(tmpPath[len(tmpPath)-1].x), float32(tmpPath[len(tmpPath)-1].y),
			float32(tmpLine[i].x), float32(tmpLine[i].y), aMap) {
			tmpPath = append(tmpPath, tmpLine[i-1])
		}
	}
	tmpPath = append(tmpPath, tmpLine[len(tmpLine)-1])

	return tmpPath
}

func (a *AStar) toFPoint(x1, y1, x2, y2 float32, bAddStart, bAddEnd bool, path []*point) (pathF []*VectorF) {
	if bAddStart {
		pathF = append(pathF, &VectorF{X: x1, Y: y1})
	}
	for i := 0; i < len(path); i++ {
		pathF = append(pathF, &VectorF{X: float32(path[i].x), Y: float32(path[i].y)})
	}
	if !bAddStart {
		pathF[0].X = x1
		pathF[0].Y = y1
	}
	if !bAddEnd {
		pathF[len(pathF)-1].X = x2
		pathF[len(pathF)-1].Y = y2
	}

	return
}

//向量单位化
func (a *AStar) vectorNormalized(x, y float32) (float32, float32) {
	a1 := x*x + y*y
	b1 := math.Sqrt(float64(a1))
	c1 := 1 / b1
	return float32(c1 * float64(x)), float32(c1 * float64(y))
}

//向量长度
func (a *AStar) vectorMagnitude(x, y float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y)))
}

//两向量间的夹角
func (a *AStar) angle(v1x, v1y, v2x, v2y float32) float32 {
	return float32(math.Acos(float64(v1x*v2x+v1y*v2y))) / (a.vectorMagnitude(v1x, v1y) * a.vectorMagnitude(v2x, v2y)) * 180.0 / math.Pi
}

func (a *AStar) approximatelyF32(a1, b1 float32) bool {
	diff := a1 - b1
	return diff > -0.0001 && diff < 0.0001
}

func (a *AStar) round(f float32) float32 {
	if f-float32(math.Floor(float64(f))) != .5 {
		return float32(math.Floor(float64(f) + .5))
	}
	r := int32(f)
	if r%2 == 0 {
		return float32(r)
	}
	return float32(r + 1)
}

func (a *AStar) clampInt(value, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

//移除一定角度内的点
func (a *AStar) removePointsOnSameLine(path []*VectorF, aMap *AMap) []*VectorF {
	pathLen := len(path)
	if pathLen <= 3 {
		return path
	}

	removeXY := make([]bool, pathLen)
	for i := 0; i < pathLen-2; {
		for j := i + 1; j < pathLen-1; {
			a1X, a1Y := a.vectorNormalized(path[i+1].X-path[i].X, path[i+1].Y-path[i].Y)
			a2X, a2Y := a.vectorNormalized(path[j+1].X-path[i].X, path[j+1].Y-path[i].Y)

			angle := a.angle(a1X, a1Y, a2X, a2Y)

			if angle <= 10 {
				//检测是否有移动阻挡
				if !a.DetectMoveCollisionBetween(path[i].X, path[i].Y, path[j+1].X, path[j+1].Y, aMap) {
					removeXY[j] = true
				} else {
					i = j + 1
					break
				}
			} else {
				i = j + 1
				break
			}

			j++
			if j == pathLen-1 {
				i = j
			}
		}
	}

	for i := len(removeXY) - 1; i > 0; i-- {
		del := removeXY[i]
		if del {
			//向前移一位
			path = append(path[:i-1], path[i:]...)
		}
	}
	return path
}

type Vectors []*VectorF

func (v Vectors) String() string {
	str := fmt.Sprintf("len: %d, ", len(v))
	for _, _v := range v {
		str += fmt.Sprintf("(%v,%v);", _v.X, _v.Y)
	}
	return str
}

func (a *AStar) smooth(path []*VectorF, aMap *AMap) []*VectorF {
	pathLen := len(path)
	if pathLen < 2 {
		return path
	}

	var maxSegmentLength float32 = 2
	var pathLength float32 = 0
	xdist := float32(0)
	ydist := float32(0)
	for i := 0; i < pathLen-1; i++ {
		xdist = path[i].X - path[i+1].X
		ydist = path[i].Y - path[i+1].Y

		pathLength += a.vectorMagnitude(xdist, ydist)
	}

	estimatedNumberOfSegments := int32(math.Floor(float64(pathLength / maxSegmentLength)))
	subdivided := make([]*VectorF, 0, estimatedNumberOfSegments+2)
	var distanceAlong float32 = 0

	for i := 0; i < pathLen-1; i++ {
		start := path[i]
		end := path[i+1]
		startX, startY := path[i].X, path[i].Y
		endX, endY := path[i+1].X, path[i+1].Y

		xdist = startX - endX
		ydist = startY - endY
		length := a.vectorMagnitude(xdist, ydist)

		for {
			if distanceAlong >= length {
				break
			}
			d := distanceAlong / length

			subdivided = append(subdivided, &VectorF{start.X + (end.X-start.X)*d, start.Y + (end.Y-start.Y)*d})
			distanceAlong += maxSegmentLength
		}

		distanceAlong -= length
	}

	subdivided = append(subdivided, path[pathLen-1])

	var iterations int = 2
	var strength float32 = 0.5

	for it := 0; it < iterations; it++ {
		prev := subdivided[0]

		for i := 1; i < len(subdivided)-1; i++ {
			tmp := subdivided[i]

			tmp2X := (prev.X + subdivided[i+1].X) / 2
			tmp2Y := (prev.Y + subdivided[i+1].Y) / 2

			subdivided[i].X = tmp.X + (tmp2X-tmp.X)*strength
			subdivided[i].Y = tmp.Y + (tmp2Y-tmp.Y)*strength

			prev = tmp
		}
	}

	return a.removePointsOnSameLine(subdivided, aMap)
}

func (a *AStar) FindPath(x1, y1, x2, y2 float32, aMap *AMap) (pathF []*VectorF, err error) {
	//检查起点，目标点是否可以寻，不能则找最近的点
	source, target, bAddStart, bAddEnd, err := a.correctPoint(x1, y1, x2, y2, aMap)
	if nil != err {
		return nil, err
	}

	//起点终点重合
	if target.x == source.x && target.y == source.y {
		pathF = append(pathF, &VectorF{X: x1, Y: y1})
		pathF = append(pathF, &VectorF{X: x2, Y: y2})
		return pathF, nil
	}

	path := a.find(source, target, aMap)
	if nil == path || len(path) < 2 {
		return nil, fmt.Errorf("no way to go.")
	}

	path = a.removePoint(path, aMap)
	pathF = a.toFPoint(x1, y1, x2, y2, bAddStart, bAddEnd, path)
	if a.bSmooth {
		pathF = a.smooth(pathF, aMap)
	}

	return
}

func (a *AStar) Print(pathF []*VectorF, aMap *AMap) {
	bHave := false
	fmt.Print("\n")
	for y := aMap.maxY - 1; y >= 0; y-- {
		for x := 0; x < aMap.maxX; x++ {
			bHave = false
			for _, point := range pathF {
				if int(point.X) == x && int(point.Y) == y {
					if !aMap.canMove(x, y) {
						fmt.Print("&")
					} else {
						fmt.Print("*")
					}
					bHave = true
					break
				}
			}
			if !bHave {
				if !aMap.canMove(x, y) {
					fmt.Print("#")
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

//起点终点不可达时搜索范围
func (a *AStar) SearchRange(iRange int) {
	a.searchRange = iRange
}

//是否启用平滑处理
func (a *AStar) UseSmooth(bSmooth bool) {
	a.bSmooth = bSmooth
}
