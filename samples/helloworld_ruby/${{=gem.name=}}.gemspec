# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require '${{=gem.name=}}/version'

Gem::Specification.new do |spec|
  spec.name         = '${{=gem.name=}}'
  spec.version      = ${{=gem.module=}}::VERSION
  spec.authors      = ['${{=gem.author=}}']
  spec.email        = ['${{=gem.author_email=}}']
  spec.summary      = %q{A sample gem}
  spec.description  = %q{Sample}
  spec.homepage     = 'http://example.com'
  spec.license      = 'All rights reserved'

  spec.files         = Dir['lib/**/*.rb'] + Dir['bin/*'] + Dir['README.md'] + Dir['docs/**/*'] - Dir['**/*~']
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.require_paths = ['lib']

  spec.required_ruby_version = '~> 2.0.0'

  spec.add_development_dependency 'bundler', '~> 1.7'
  spec.add_development_dependency 'rspec', '= 3.1.0'
end
