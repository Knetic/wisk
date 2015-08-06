#!/usr/bin/env python
# coding=utf-8

from setuptools import setup, find_packages


setup(name='${{=project.name=}}',
      version='0.0.1',
      author='${{=project.author=}}',
      author_email='${{=project.author_email=}}',
      url='http://www.${{=project.url=}}',
      download_url='http://www.${{=project.name=}}/files/',
      description='Short description of ${{=project.name=}}...',
      long_description='Short description of ${{=project.name=}}...',

      packages = find_packages(),
      include_package_data = True,
      package_data = {
        '': ['*.txt', '*.rst'],
        '${{=project.name=}}': ['data/*.html', 'data/*.css'],
      },
      exclude_package_data = { '': ['README.txt'] },

      scripts = ['bin/${{=project.name=}}'],

      keywords='python tools utils internet www',
      license='GPL',
      classifiers=['Development Status :: 5 - Production/Stable',
                   'Natural Language :: English',
                   'Operating System :: OS Independent',
                   'Programming Language :: Python :: 2',
                   'License :: OSI Approved :: GNU Library or Lesser General Public License (LGPL)',
                   'License :: OSI Approved :: GNU Affero General Public License v3',
                   'Topic :: Internet',
                   'Topic :: Internet :: WWW/HTTP',
                  ],

      #setup_requires = ['python-stdeb', 'fakeroot', 'python-all'],
      install_requires = ['setuptools'],
     )
