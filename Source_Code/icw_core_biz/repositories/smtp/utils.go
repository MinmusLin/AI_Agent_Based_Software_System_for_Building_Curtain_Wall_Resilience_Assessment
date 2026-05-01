package smtp

import (
	"errors"
	"fmt"
	"html"

	"icw_core_biz/internal/services/auth/consts"
)

// emailSceneText 根据邮箱验证码业务场景类型生成邮件标题和业务场景名称
func emailSceneText(scene string) (string, string, error) {
	switch consts.ParseEmailCodeScene(scene) {
	case consts.SceneRegister:
		return "注册验证码 - 建筑幕墙韧性评估软件系统", "账号注册", nil
	case consts.SceneLogin:
		return "登录验证码 - 建筑幕墙韧性评估软件系统", "账号登录", nil
	case consts.SceneReset:
		return "重置验证码 - 建筑幕墙韧性评估软件系统", "重置密码", nil
	default:
		return "", "", errors.New("invalid email code scene")
	}
}

// buildEmailCodeHTML 构建验证码邮件 HTML 正文
func buildEmailCodeHTML(sceneName, code string, ttlMinutes int) string {
	return fmt.Sprintf(`
		<!doctype html>
		<html>
		<body style="margin:0;padding:0;background:#f6f8fb;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Arial,sans-serif;color:#172033;">
			<div style="max-width:560px;margin:0 auto;padding:32px 16px;">
				<div style="background:#ffffff;border:1px solid #e7ebf0;border-radius:8px;padding:28px;">
					<div style="font-size:18px;font-weight:700;margin-bottom:8px;">建筑幕墙韧性评估软件系统</div>
					<div style="font-size:14px;color:#5f6b7a;margin-bottom:24px;">%s验证码</div>
					<div style="font-size:13px;color:#5f6b7a;margin-bottom:10px;">您正在进行%s操作，验证码为：</div>
					<div style="letter-spacing:8px;font-size:32px;font-weight:700;color:#1565c0;background:#f1f6ff;border-radius:6px;padding:16px 18px;text-align:center;">%s</div>
					<div style="font-size:13px;line-height:1.7;color:#5f6b7a;margin-top:24px;">该验证码 %d 分钟内有效，请勿将验证码告知任何人。若非您本人操作，请忽略本邮件。</div>
					<div style="font-size:12px;line-height:1.6;color:#8a95a3;margin-top:12px;font-style:italic;">邮件由系统自动发送，请勿直接回复。</div>
				</div>
			</div>
		</body>
		</html>
	`, html.EscapeString(sceneName), html.EscapeString(sceneName), html.EscapeString(code), ttlMinutes)
}
