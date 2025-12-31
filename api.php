<?php
function isMobile(): bool {
    $userAgent = $_SERVER['HTTP_USER_AGENT'] ?? '';

    $pattern = '/Android|iP(hone|od)|Windows ?CE|Symbian|Mobile|Opera Mobi|BlackBerry|Palm(OS)?|PocketPC|SonyEricsson|Nokia|Vodafone|HTC_|SAMSUNG-SGH|armv[56]l|Go\.Web|J2ME\/MIDP/i';

    return preg_match($pattern, $userAgent) > 0;
}

$pcPage = 'pc.php';
$mobilePage = 'pe.php';

header('Location: ' . (isMobile() ? $mobilePage : $pcPage));
exit;