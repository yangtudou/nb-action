package bark

import (
	"testing"
)

func TestEncodeCiphertext(t *testing.T) {

	// 测试输入
	data := []byte("hello bark")

	// 执行编码
	result := encodeCiphertext(data)

	// 期望结果
	expected := "aGVsbG8gYmFyaw%3D%3D"

	// 对比结果

	if result != expected {

		t.Fatalf(
			"unexpected result\nwant: %s\ngot: %s",
			expected,
			result,
		)

	}

	t.Log("encode success")

}
