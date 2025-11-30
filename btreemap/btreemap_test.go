package btreemap

import (
	"testing"
)

// 整数比较函数
func intComparator(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// TestBTreeMapBasic 测试基本的键值对操作
func TestBTreeMapBasic(t *testing.T) {
	// 创建一个B树映射，阶数为3
	m := New[int, string](intComparator)

	// 测试Put和Get操作
	m.Insert(3, "three")
	m.Insert(1, "one")
	m.Insert(4, "four")

	// 验证大小
	if m.Len() != 3 {
		t.Errorf("Expected size 3, got %d", m.Len())
	}

	// 验证Get操作
	val, found := m.Get(1)
	if !found || val != "one" {
		t.Errorf("Expected 'one', got %v, found: %v", val, found)
	}

	val, found = m.Get(3)
	if !found || val != "three" {
		t.Errorf("Expected 'three', got %v, found: %v", val, found)
	}

	val, found = m.Get(4)
	if !found || val != "four" {
		t.Errorf("Expected 'four', got %v, found: %v", val, found)
	}

	// 测试不存在的键
	_, found = m.Get(10)
	if found {
		t.Error("Expected not found for key 10")
	}

	// 测试ContainsKey
	if !m.ContainsKey(3) {
		t.Error("Expected ContainsKey(3) to return true")
	}

	if m.ContainsKey(10) {
		t.Error("Expected ContainsKey(10) to return false")
	}

	// 测试更新已存在的键
	m.Insert(3, "THREE")
	val, found = m.Get(3)
	if !found || val != "THREE" {
		t.Errorf("Expected 'THREE' after update, got %v", val)
	}
	if m.Len() != 3 {
		t.Errorf("Expected size to remain 3 after update, got %d", m.Len())
	}

	// 测试删除键
	m.Remove(1)
	if m.Len() != 2 {
		t.Errorf("Expected size 2 after deletion, got %d", m.Len())
	}

	_, found = m.Get(1)
	if found {
		t.Error("Expected not found for deleted key 1")
	}

	// 测试删除不存在的键
	m.Remove(10)
	if m.Len() != 2 {
		t.Errorf("Expected size to remain 2 after deleting non-existent key, got %d", m.Len())
	}
}

// TestBTreeMapClear 测试Clear操作
func TestBTreeMapClear(t *testing.T) {
	m := New[int, string](intComparator)

	// 添加一些键值对
	for i := 0; i < 10; i++ {
		m.Insert(i, string(rune(i+'a')))
	}

	// 验证大小
	if m.Len() != 10 {
		t.Errorf("Expected size 10, got %d", m.Len())
	}

	// 清除映射
	m.Clear()

	// 验证映射是否为空
	if !m.IsEmpty() {
		t.Error("Map should be empty after Clear()")
	}

	if m.Len() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", m.Len())
	}
}

// 辅助函数：收集键到切片
func collectKeys(keysIter func(func(int) bool)) []int {
	var keys []int
	keysIter(func(k int) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

// 辅助函数：收集值到切片
func collectValues(valuesIter func(func(string) bool)) []string {
	var values []string
	valuesIter(func(v string) bool {
		values = append(values, v)
		return true
	})
	return values
}

// 测试node结构体和相关方法
func TestBTreeMapEntryAndEntries(t *testing.T) {
	// 测试1: 直接测试node结构体
	entry := node[int, string]{Key: 42, Value: "answer"}
	if entry.GetKey() != 42 || entry.GetValue() != "answer" {
		t.Errorf("node creation failed, expected (42, 'answer'), got (%v, %v)", entry.GetKey(), entry.GetValue())
	}

	// 测试2: 测试Entries方法
	m := New[int, string](intComparator)

	// 添加键值对
	m.Insert(3, "three")
	m.Insert(1, "one")
	m.Insert(2, "two")

	// 使用Range方法收集条目并构建node结构体
	var collectedEntries []node[int, string]
	for k, v := range m.Range(nil, nil) {
		collectedEntries = append(collectedEntries, node[int, string]{Key: k, Value: v})
	}

	// 验证条目数量
	if len(collectedEntries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(collectedEntries))
	}

	// 验证条目内容，应该按键的升序排列
	expectedEntries := []node[int, string]{
		{Key: 1, Value: "one"},
		{Key: 2, Value: "two"},
		{Key: 3, Value: "three"},
	}

	for i, e := range collectedEntries {
		expected := expectedEntries[i]
		if e.Key != expected.Key || e.Value != expected.Value {
			t.Errorf("Expected entry (%v, %v) at position %d, got (%v, %v)",
				expected.Key, expected.Value, i, e.Key, e.Value)
		}

		// 测试Entry的Getter方法
		if e.GetKey() != e.Key || e.GetValue() != e.Value {
			t.Errorf("Entry Getter methods failed for entry (%v, %v)", e.Key, e.Value)
		}
	}

	// 测试3: 测试空映射的Range方法
	emptyMap := New[int, string](intComparator)
	emptyEntriesCount := 0
	for range emptyMap.Range(nil, nil) {
		emptyEntriesCount++
	}
	if emptyEntriesCount != 0 {
		t.Errorf("Expected 0 entries for empty map, got %d", emptyEntriesCount)
	}
}

// 测试Range方法的范围查询功能
func TestBTreeMapRangeEntries(t *testing.T) {
	m := New[int, string](intComparator)

	// 添加键值对
	for i := 1; i <= 10; i++ {
		m.Insert(i, "value-"+string(rune('0'+i)))
	}

	// 测试4: 测试Range方法，带上下界
	lowerBound := 3
	upperBound := 7
	var rangeEntries []node[int, string]
	for k, v := range m.Range(&lowerBound, &upperBound) {
		rangeEntries = append(rangeEntries, node[int, string]{Key: k, Value: v})
	}

	// 应该返回键为3,4,5,6的条目（不包含7）
	expectedKeys := []int{3, 4, 5, 6}
	if len(rangeEntries) != len(expectedKeys) {
		t.Errorf("Expected %d entries in range, got %d", len(expectedKeys), len(rangeEntries))
	}

	for i, e := range rangeEntries {
		if e.Key != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, e.Key)
		}
		if e.Value != "value-"+string(rune('0'+e.Key)) {
			t.Errorf("Expected value 'value-%d' for key %d, got '%s'", e.Key, e.Key, e.Value)
		}
	}

	// 测试5: 测试Range方法，只有下界
	lowerOnly := 8
	var lowerRangeEntries []node[int, string]
	for k, v := range m.Range(&lowerOnly, nil) {
		lowerRangeEntries = append(lowerRangeEntries, node[int, string]{Key: k, Value: v})
	}

	// 应该返回键为8,9,10的条目
	expectedLowerKeys := []int{8, 9, 10}
	if len(lowerRangeEntries) != len(expectedLowerKeys) {
		t.Errorf("Expected %d entries with lower bound, got %d", len(expectedLowerKeys), len(lowerRangeEntries))
	}

	// 测试6: 测试Range方法，只有上界
	upperOnly := 2
	var upperRangeEntries []node[int, string]
	for k, v := range m.Range(nil, &upperOnly) {
		upperRangeEntries = append(upperRangeEntries, node[int, string]{Key: k, Value: v})
	}

	// 应该返回键为1的条目（不包含2）
	expectedUpperKeys := []int{1}
	if len(upperRangeEntries) != len(expectedUpperKeys) {
		t.Errorf("Expected %d entries with upper bound, got %d", len(expectedUpperKeys), len(upperRangeEntries))
	}
}

func TestBTreeMapEntryAPI(t *testing.T) {
	// 创建一个BTreeMap实例
	m := New[int, string](intComparator)

	// 测试1: 测试不存在的键的Entry状态 (通过ContainsKey检查)
	entry1 := m.Entry(10)
	if m.ContainsKey(10) {
		t.Errorf("Expected key 10 to not exist")
	}

	// 测试2: OrInsert - 不存在的键
	valPtr1 := entry1.OrInsert("ten")
	if *valPtr1 != "ten" {
		t.Errorf("Expected OrInsert to return 'ten', got '%s'", *valPtr1)
	}

	// 验证值被正确插入
	value, found := m.Get(10)
	if !found || value != "ten" {
		t.Errorf("Expected value 'ten' to be inserted, got %v, found: %v", value, found)
	}

	// 测试3: 测试已存在的键的Entry状态 (通过ContainsKey检查)
	entry2 := m.Entry(10)
	if !m.ContainsKey(10) {
		t.Errorf("Expected key 10 to exist")
	}

	// 测试4: OrInsert - 已存在的键
	valPtr2 := entry2.OrInsert("TEN")
	if *valPtr2 != "ten" {
		t.Errorf("Expected OrInsert to return existing value 'ten', got '%s'", *valPtr2)
	}

	// 验证值没有被覆盖
	value, found = m.Get(10)
	if !found || value != "ten" {
		t.Errorf("Expected value 'ten' to remain unchanged, got %v, found: %v", value, found)
	}

	// 测试5: OrInsertWith - 不存在的键
	counter := 0
	entry3 := m.Entry(20)
	valPtr3 := entry3.OrInsertWith(func() string {
		counter++
		return "twenty"
	})

	if *valPtr3 != "twenty" {
		t.Errorf("Expected OrInsertWith to return 'twenty', got '%s'", *valPtr3)
	}
	if counter != 1 {
		t.Errorf("Expected OrInsertWith to call function once, called %d times", counter)
	}

	// 测试6: OrInsertWith - 已存在的键
	entry4 := m.Entry(20)
	valPtr4 := entry4.OrInsertWith(func() string {
		counter++
		return "TWENTY"
	})

	if *valPtr4 != "twenty" {
		t.Errorf("Expected OrInsertWith to return existing value 'twenty', got '%s'", *valPtr4)
	}
	if counter != 1 {
		t.Errorf("Expected OrInsertWith to not call function for existing key, called %d times", counter)
	}

	// 测试7: Get方法 - 已存在的键
	val5, found5 := entry4.Get()
	if !found5 || val5 != "twenty" {
		t.Errorf("Expected Get to return 'twenty' and true, got '%v', found: %v", val5, found5)
	}

	// 测试8: Get方法 - 不存在的键
	entry5 := m.Entry(30)
	val6, found6 := entry5.Get()
	if found6 {
		t.Errorf("Expected Get to return false for non-existent key, got found: %v, value: '%v'", found6, val6)
	}

	// 测试9: Insert方法 - 不存在的键
	valPtr7 := entry5.Insert("thirty")
	if *valPtr7 != "thirty" {
		t.Errorf("Expected Insert to return 'thirty', got '%s'", *valPtr7)
	}

	// 测试10: Insert方法 - 已存在的键
	valPtr8 := entry4.Insert("TWENTY")
	if *valPtr8 != "TWENTY" {
		t.Errorf("Expected Insert to return 'TWENTY', got '%s'", *valPtr8)
	}

	// 验证值被覆盖
	value, found = m.Get(20)
	if !found || value != "TWENTY" {
		t.Errorf("Expected value 'TWENTY' to be inserted, got %v, found: %v", value, found)
	}

	// 测试11: 使用Insert方法修改值
	entry6 := m.Entry(10)
	entry6.Insert("TEN_MODIFIED")

	value, found = m.Get(10)
	if !found || value != "TEN_MODIFIED" {
		t.Errorf("Expected value to be modified to 'TEN_MODIFIED', got %v, found: %v", value, found)
	}
}

// TestBTreeMapForEach 测试ForEach方法
func TestBTreeMapForEach(t *testing.T) {
	m := New[int, string](intComparator)

	// 添加键值对
	m.Insert(3, "three")
	m.Insert(1, "one")
	m.Insert(2, "two")

	// 用于存储遍历结果
	var keys []int
	var values []string

	// 使用Iter遍历替代ForEach
	for k, v := range m.Iter() {
		keys = append(keys, k)
		values = append(values, v)
	}

	// 验证遍历顺序
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys in ForEach, got %d", len(expectedKeys), len(keys))
	}

	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}

	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values in ForEach, got %d", len(expectedValues), len(values))
	}

	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %s at position %d, got %s", expectedValues[i], i, v)
		}
	}
}

// TestBTreeMapRange 测试Range方法的迭代器功能
func TestBTreeMapRange(t *testing.T) {
	m := New[int, string](intComparator)

	// 添加键值对
	for i := 1; i <= 5; i++ {
		m.Insert(i, string(rune(i+'a'-1)))
	}

	// 用于存储遍历结果
	var keys []int
	var values []string

	// 使用Iter迭代器遍历所有元素
	for k, v := range m.Iter() {
		keys = append(keys, k)
		values = append(values, v)
	}

	// 验证遍历结果，应该包含所有键值对且按升序排列
	expectedKeys := []int{1, 2, 3, 4, 5}
	expectedValues := []string{"a", "b", "c", "d", "e"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys in Range iteration, got %d", len(expectedKeys), len(keys))
	}

	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}

	if len(values) != len(expectedValues) {
		t.Errorf("Expected %d values in Range iteration, got %d", len(expectedValues), len(values))
	}

	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Expected value %s at position %d, got %s", expectedValues[i], i, v)
		}
	}

	// 测试提前终止遍历
	keys = nil
	values = nil
	for k, v := range m.Iter() {
		keys = append(keys, k)
		values = append(values, v)
		if k >= 3 {
			break // 提前终止遍历
		}
	}

	// 验证提前终止的结果
	expectedKeys = []int{1, 2, 3}
	expectedValues = []string{"a", "b", "c"}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys in Range with early termination, got %d", len(expectedKeys), len(keys))
	}

	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}
}

// TestBTreeMapLargeDataSet 测试大数据集
func TestBTreeMapLargeDataSet(t *testing.T) {
	m := New[int, int](intComparator)

	// 添加1000个键值对
	for i := 0; i < 1000; i++ {
		m.Insert(i, i*10)
	}

	// 验证大小
	if m.Len() != 1000 {
		t.Errorf("Expected size 1000, got %d", m.Len())
	}

	// 验证一些键值对
	for i := 0; i < 1000; i += 100 {
		val, found := m.Get(i)
		if !found || val != i*10 {
			t.Errorf("Expected %d for key %d, got %v, found: %v", i*10, i, val, found)
		}
	}

	// 删除一些键值对
	for i := 0; i < 1000; i += 2 {
		m.Remove(i)
	}

	// 验证剩余大小
	if m.Len() != 500 {
		t.Errorf("Expected size 500 after deleting even keys, got %d", m.Len())
	}

	// 验证删除后的键值对
	for i := 0; i < 1000; i++ {
		_, found := m.Get(i)
		if i%2 == 0 && found {
			t.Errorf("Expected key %d to be deleted", i)
		} else if i%2 == 1 && !found {
			t.Errorf("Expected key %d to exist", i)
		}
	}
}

// 辅助函数：收集字符串键到切片
func collectStringKeys(keysIter func(func(string) bool)) []string {
	var keys []string
	keysIter(func(k string) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}
