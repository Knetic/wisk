#!/usr/bin/env python
# coding=utf-8

import unittest

from ${{=project.name=}} import ${{=class.name=}}


class Test(unittest.TestCase):
    """Unit tests for ${{=project.name=}}.${{=class.name=}}"""

    def test_fn(self):
        """Test result"""
        instance = ${{=class.name=}}()
        value = True
        result = instance.exists
        self.assertEqual(value, result)

if __name__ == "__main__":
    unittest.main()
