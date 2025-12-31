<?php

// 图片URL文件
$txtFile = 'acg_横图.txt';

$lines = array_map('trim', file($txtFile));

$randomIndex = rand(0, count($lines) - 1);
$imageUrl = cleanUrl($lines[$randomIndex]);

$result = [
    "code" => "200",
    "acgurl" => $imageUrl
];

$type = $_GET['return'] ?? '';

switch ($type) {

    case 'json':
        $imageInfo = @getimagesize($imageUrl);
        if ($imageInfo) {
            $result['width'] = strval($imageInfo[0]);
            $result['height'] = strval($imageInfo[1]);
        } else {
            $result['width'] = "0";
            $result['height'] = "0";
        }

        $pathinfo = pathinfo($imageUrl);
        $result['size'] = $pathinfo['extension'] ?? '';

        header('Content-Type: application/json; charset=utf-8');
        echo json_encode($result);
        break;

    case 'all':
        // 直接返回全部 URL 内容
        header('Content-Type: text/plain; charset=utf-8');
        echo implode("\n", $lines);
        break;

    default:
        header("Location: " . $imageUrl);
        break;
}

function cleanUrl(string $str): string
{
    return str_replace([" ", "\n", "\r", "\t"], '', $str);
}
?>