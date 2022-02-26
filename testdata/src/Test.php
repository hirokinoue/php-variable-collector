<?php

namespace Test;

use Test\Baz\Bar;
use Test\Baz\Foo;

class Test 
{
    private $bar;
    private $foo;

    public function __constructor(Bar $bar, Foo $foo): void
    {
        $this->bar = $bar;
        $this->foo = $foo;
    }

    public function getBar(): Bar
    {
        return $this->bar;
    }

    public function getFoo(): Foo
    {
        return $this->foo;
    }
}