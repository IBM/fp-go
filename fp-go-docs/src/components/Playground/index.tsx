import React, { useState, useEffect, useRef } from 'react';
import type { ReactElement } from 'react';
import styles from './styles.module.css';

interface Example {
  id: string;
  title: string;
  code: string;
  output?: string;
}

const examples: Example[] = [
  {
    id: 'option',
    title: 'Option · safe values',
    code: `package main

import (
	"fmt"
	O "github.com/IBM/fp-go/v2/option"
)

func main() {
	some := O.Some(42)
	none := O.None[int]()
	
	fmt.Println("Some:", some)
	fmt.Println("None:", none)
	
	// Map over option
	doubled := O.Map(func(x int) int { return x * 2 })(some)
	fmt.Println("Doubled:", doubled)
}`,
    output: 'Some: Some(42)\nNone: None\nDoubled: Some(84)'
  },
  {
    id: 'either',
    title: 'Either · error handling',
    code: `package main

import (
	"fmt"
	E "github.com/IBM/fp-go/v2/either"
)

func divide(a, b int) E.Either[string, int] {
	if b == 0 {
		return E.Left[int]("division by zero")
	}
	return E.Right[string](a / b)
}

func main() {
	result1 := divide(10, 2)
	result2 := divide(10, 0)
	
	fmt.Println("10 / 2 =", result1)
	fmt.Println("10 / 0 =", result2)
}`,
    output: '10 / 2 = Right(5)\n10 / 0 = Left(division by zero)'
  },
  {
    id: 'array',
    title: 'Array · functional operations',
    code: `package main

import (
 "fmt"
 A "github.com/IBM/fp-go/v2/array"
)

func main() {
 numbers := []int{1, 2, 3, 4, 5}
 
 // Map: double each number
 doubled := A.Map(func(x int) int {
  return x * 2
 })(numbers)
 
 // Filter: keep only even numbers
 evens := A.Filter(func(x int) bool {
  return x%2 == 0
 })(doubled)
 
 fmt.Println("Doubled:", doubled)
 fmt.Println("Evens:", evens)
}`,
    output: 'Doubled: [2 4 6 8 10]\nEvens: [2 4 6 8 10]'
  },
  {
    id: 'pipe',
    title: 'Pipe · function composition',
    code: `package main

import (
 "fmt"
 F "github.com/IBM/fp-go/v2/function"
)

func main() {
 addOne := func(x int) int { return x + 1 }
 double := func(x int) int { return x * 2 }
 square := func(x int) int { return x * x }
 
 // Compose functions: (5 + 1) * 2 = 12, then 12^2 = 144
 result := F.Pipe3(
  5,
  addOne,
  double,
  square,
 )
 
 fmt.Println("Result:", result)
}`,
    output: 'Result: 144'
  },
  {
    id: 'chain',
    title: 'Chain · monadic operations',
    code: `package main

import (
 "fmt"
 O "github.com/IBM/fp-go/v2/option"
)

func safeDivide(a, b int) O.Option[int] {
 if b == 0 {
  return O.None[int]()
 }
 return O.Some(a / b)
}

func main() {
 // Chain operations together
 result := O.Chain(func(x int) O.Option[int] {
  return safeDivide(x, 2)
 })(O.Some(20))
 
 fmt.Println("20 / 2 =", result)
 
 // Chain with None
 none := O.Chain(func(x int) O.Option[int] {
  return safeDivide(x, 0)
 })(O.Some(20))
 
 fmt.Println("20 / 0 =", none)
}`,
    output: '20 / 2 = Some(10)\n20 / 0 = None'
  },
  {
    id: 'fold',
    title: 'Fold · reduce arrays',
    code: `package main

import (
 "fmt"
 A "github.com/IBM/fp-go/v2/array"
)

func main() {
 numbers := []int{1, 2, 3, 4, 5}
 
 // Sum all numbers
 sum := A.Reduce(func(acc, x int) int {
  return acc + x
 }, 0)(numbers)
 
 // Product of all numbers
 product := A.Reduce(func(acc, x int) int {
  return acc * x
 }, 1)(numbers)
 
 fmt.Println("Sum:", sum)
 fmt.Println("Product:", product)
}`,
    output: 'Sum: 15\nProduct: 120'
  }
];

export default function Playground(): ReactElement {
  const [activeExample, setActiveExample] = useState(0);
  const [code, setCode] = useState(examples[0].code);
  const [output, setOutput] = useState(examples[0].output || '');
  const [isRunning, setIsRunning] = useState(false);

  const handleExampleChange = (index: number) => {
    setActiveExample(index);
    setCode(examples[index].code);
    setOutput(examples[index].output || '');
  };

  const handleRun = async () => {
    setIsRunning(true);
    setOutput('Running...');
    
    try {
      // Use production fpgo-sandbox endpoint
      const response = await fetch('https://fpgo-sandbox.fly.dev/v1/exec', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          sandbox: 'go',
          command: 'run',
          files: {
            '': code  // Empty string key as per API spec
          }
        })
      });

      const result = await response.json();
      
      if (result.ok) {
        setOutput(result.stdout || result.output || 'Success!');
      } else {
        setOutput(`Error: ${result.stderr || result.error || 'Execution failed'}`);
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to execute code';
      setOutput(`Error: ${errorMessage}`);
    } finally {
      setIsRunning(false);
    }
  };

  const lineCount = code.split('\n').length;

  return (
    <div className={styles.page}>
      <div className={styles.eyebrow}>Playground</div>
      <h1 className={styles.title}>
        Try fp-go in your <em>browser.</em>
      </h1>
      <p className={styles.sub}>
        Edit, run and share Go snippets without leaving the docs. Compiled & executed via codapi — no installs, no setup.
      </p>

      <div className={styles.snippets}>
        {examples.map((example, index) => (
          <button
            key={example.id}
            className={`${styles.snippet} ${index === activeExample ? styles.active : ''}`}
            onClick={() => handleExampleChange(index)}
          >
            <span className={styles.ico}>{index === activeExample ? '●' : '○'}</span>
            {example.title}
          </button>
        ))}
      </div>

      <div className={styles.pg}>
        {/* Toolbar */}
        <div className={styles.pgToolbar}>
          <div className={styles.pgTabs}>
            <button className={`${styles.pgTab} ${styles.active}`}>
              <span className={styles.dot}></span>
              main.go
            </button>
          </div>
          <div className={styles.spacer}></div>
          <div className={styles.pgMeta}>
            <span className={styles.lbl}>Runtime</span>
            <span className={styles.pill}>go 1.22</span>
            <span className={styles.pill}>fp-go v2.2.82</span>
          </div>
        </div>

        {/* Body */}
        <div className={styles.pgBody}>
          <div className={styles.editor}>
            <div className={styles.gutter}>
              {Array.from({ length: lineCount }, (_, i) => (
                <span key={i + 1}>{i + 1}</span>
              ))}
            </div>
            <textarea
              className={styles.codeArea}
              value={code}
              onChange={(e) => setCode(e.target.value)}
              spellCheck={false}
            />
          </div>

          <aside className={styles.side}>
            <div className={styles.sideTabs}>
              <button className={`${styles.sideTab} ${styles.active}`}>
                Output
              </button>
            </div>

            <div className={styles.stdout}>
              <span className={styles.ts}>// stdout</span>
              <div className={styles.val}>{output}</div>
            </div>
          </aside>
        </div>

        {/* Action bar */}
        <div className={styles.pgActionbar}>
          <button 
            className={`${styles.btn} ${styles.primary}`}
            onClick={handleRun}
            disabled={isRunning}
          >
            <svg width="14" height="14" viewBox="0 0 32 32" fill="currentColor">
              <path d="M11 23v-14l11 7-11 7z"/>
            </svg>
            {isRunning ? 'Running...' : 'Run'}
            <span className={styles.kbd}>⌘↵</span>
          </button>
          <div className={styles.spacer} />
          <div className={styles.by}>
            <span className={styles.done}>
              <svg width="12" height="12" viewBox="0 0 32 32" fill="currentColor">
                <path d="M13 24l-9-9 1.4-1.4L13 21.2 26.6 7.6 28 9z"/>
              </svg>
              Ready
            </span>
            <span style={{color: 'var(--gray-50)'}}>·</span>
            <span>powered by <a href="https://codapi.org">codapi</a></span>
          </div>
        </div>

        {/* Status bar */}
        <div className={styles.pgStatus}>
          <span className={styles.seg}>
            <span className={styles.dot}></span> Connected
          </span>
          <span className={styles.seg}>
            <span className={styles.k}>go</span>
            <span className={styles.v}>1.22.1</span>
          </span>
          <div className={styles.spacer} />
          <span className={styles.seg}>
            <span className={styles.k}>UTF-8</span>
          </span>
        </div>
      </div>
    </div>
  );
}

// Made with Bob
