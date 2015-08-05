require '${{=gem.name=}}/version'

${{=:gem.module=}}
module ${{=_value=}}
  ${{=_recurse=}}
end
${{=;gem.module=}}

${{=gem.module[::]=}}.VERSION = '1.0.0'
