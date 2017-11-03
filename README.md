Finds shortest paths from node with vertex ID 1 to all other vertices in graph.

Takes input as file with vertices identified in this structure:

~~~
[vertex ID] [head vertex,edge weight] [head vertex,edge weight]
~~~

For example:

~~~
1 2,100 3,200
2 1,100 3,150
3 1,200 2,150
~~~

Output is randomly ordered, but pairs are identical.