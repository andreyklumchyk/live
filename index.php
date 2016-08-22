<?php

$ch = curl_init();
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$url = 'http://image.slidesharecdn.com/random-160304224949/95/-';

for ($i = 1; $i < 408; $i++) {
    $is_loaded = false;
    foreach ([1024, 638, 320] as $size) {
        if ($is_loaded) {
            continue;
        }
        $file_name = $i.'-'.$size.'.jpg';
        $file = $url.$file_name;
        curl_setopt($ch, CURLOPT_URL, $file);
        $output = curl_exec($ch);
        if (strpos(substr($output, 0, 100), 'AccessDenied') !== false) {
            continue;
        }
        file_put_contents('result/'.$i.'.jpg', $output);
        $is_loaded = true;
    }
    if (!$is_loaded) {
        print_r('File missed: '.$i);
    }
}

curl_close($ch);
