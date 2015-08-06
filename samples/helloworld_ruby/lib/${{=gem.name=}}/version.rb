require '${{=gem.name=}}/version'

${{=:gem.module=}}
module ${{value}}
  ${{recurse}}
end
${{=;gem.module=}}

${{=gem.module[::]=}}.VERSION = '1.0.0'
