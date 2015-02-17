gostorages
==========

An unified interface to manipulate storage engine for Go.

gostorages is used in `picfit <https://github.com/thoas/picfit>`_ to allow us
switching over storage engine.

Currently, it supports the following storages:

* Amazon S3
* File system

Installation
============

Just run:

::
    $ go get github.com/thoas/gostorages

Usage
=====

File system
-----------

.. code-block:: go

    package main

    import (
        "fmt"
        "github.com/thoas/gostorages
        "os"
    )

	tmp := os.TempDir()

	storage := gostorages.NewFileSystemStorage(tmp, "http://img.example.com")

    // Saving a file named test
    storage.Save("test", gostorages.ContentFile([]byte("(╯°□°）╯︵ ┻━┻")))

    storage.URL("test") // => http://img.example.com/test

    // Deleting the new file on the storage
    storage.Delete("test")

Roadmap
=======

see `issues <https://github.com/thoas/gostorages/issues>`_

Don't hesitate to send patch or improvements.
