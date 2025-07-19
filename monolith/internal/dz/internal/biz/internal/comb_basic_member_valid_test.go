package biz_internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TODO
// код теста не мой! Взял его до поиска более разумного решения как отследить запуск всех функций в сборных функциях: например если кто-то что-то закоментит временно и забудет раскомментировать. Потому что в таком случае юнит-тесты покажут что все норм, а проверки фактически не будет! Поэтому нельзя такого допускать. Но мне кажется должно быть более простое решение, не думаю что разработчики go не предусмотрели такую очевидную потребность

func TestBasicMemberValidation_AllChecksCalled_AST(t *testing.T) {
	// 1. Настраиваем парсер для текущего файла
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "comb_basic_member_valid.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("Ошибка парсинга файла comb_basic_member_valid.go: %v", err)
	}

	// 2. Ищем метод BasicMemberValidation в AST
	var method *ast.FuncDecl
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == "BasicMemberValidation" {
			method = fd
			break
		}
	}
	if method == nil {
		t.Fatal("Метод BasicMemberValidation не найден")
	}

	// 3. Ожидаемые проверки
	requiredChecks := []string{
		"b.UuidCheck",
		"b.nameCheck",
		"b.phoneCheck",
		"b.communityIdCheck",
	}

	// 4. Собираем фактические вызовы в методе
	foundChecks := make(map[string]bool)
	ast.Inspect(method.Body, func(n ast.Node) bool {
		// Ищем только вызовы функций (CallExpr)
		if call, ok := n.(*ast.CallExpr); ok {
			// Проверяем, что это вызов метода через b. (SelectorExpr)
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				// Проверяем, что вызывается метод структуры (b.xxx)
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "b" {
					checkName := "b." + sel.Sel.Name
					foundChecks[checkName] = true
				}
			}
		}
		return true
	})

	// 5. Проверяем наличие всех обязательных проверок
	for _, check := range requiredChecks {
		if !foundChecks[check] {
			t.Errorf("Проверка %s отсутствует в методе (возможно закомментирована)", check)
		}
	}

	// 6. Дополнительная проверка: нет ли лишних проверок?
	if len(foundChecks) != len(requiredChecks) {
		t.Errorf(
			"Несоответствие количества проверок: ожидалось %d, найдено %d",
			len(requiredChecks),
			len(foundChecks),
		)
	}
}
