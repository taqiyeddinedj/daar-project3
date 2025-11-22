#!/usr/bin/env python3
"""
DAAR Project 3 - Performance Chart Generator
Generates bar charts with error bars for the LaTeX report
"""

import json
import matplotlib.pyplot as plt
import numpy as np

# Set style for academic papers
plt.style.use('seaborn-v0_8-darkgrid')
plt.rcParams['figure.figsize'] = (10, 6)
plt.rcParams['font.size'] = 11
plt.rcParams['axes.labelsize'] = 12
plt.rcParams['axes.titlesize'] = 14
plt.rcParams['xtick.labelsize'] = 10
plt.rcParams['ytick.labelsize'] = 10

def load_results():
    """Load benchmark results from JSON file"""
    try:
        with open('benchmark_results.json', 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        print("Error: benchmark_results.json not found!")
        print("Run: go run benchmark_test.go first")
        return None

def plot_simple_search(data):
    """Bar chart for simple search with error bars"""
    queries = [item['query'] for item in data]
    means = [item['mean_ms'] for item in data]
    stddevs = [item['stddev_ms'] for item in data]
    
    fig, ax = plt.subplots(figsize=(10, 6))
    
    x = np.arange(len(queries))
    bars = ax.bar(x, means, yerr=stddevs, capsize=5, 
                   color='#2ecc71', edgecolor='black', linewidth=1.2,
                   error_kw={'linewidth': 2, 'ecolor': '#e74c3c'})
    
    ax.set_xlabel('Requête de recherche', fontweight='bold')
    ax.set_ylabel('Temps (ms)', fontweight='bold')
    ax.set_title('Performance de la recherche simple\n(Moyenne ± Écart-type, n=100)', 
                 fontweight='bold', pad=20)
    ax.set_xticks(x)
    ax.set_xticklabels(queries, rotation=0)
    ax.grid(axis='y', alpha=0.3)
    
    # Add value labels on bars
    for i, (bar, mean, std) in enumerate(zip(bars, means, stddevs)):
        height = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2., height + std + 0.1,
                f'{mean:.2f}±{std:.2f}',
                ha='center', va='bottom', fontsize=9, fontweight='bold')
    
    plt.tight_layout()
    plt.savefig('chart_simple_search.pdf', dpi=300, bbox_inches='tight')
    plt.savefig('chart_simple_search.png', dpi=300, bbox_inches='tight')
    print("✓ Generated: chart_simple_search.pdf/png")
    plt.close()

def plot_regex_search(data):
    """Bar chart for regex search with error bars"""
    queries = [item['query'] for item in data]
    means = [item['mean_ms'] for item in data]
    stddevs = [item['stddev_ms'] for item in data]
    
    fig, ax = plt.subplots(figsize=(12, 6))
    
    x = np.arange(len(queries))
    bars = ax.bar(x, means, yerr=stddevs, capsize=5,
                   color='#3498db', edgecolor='black', linewidth=1.2,
                   error_kw={'linewidth': 2, 'ecolor': '#e74c3c'})
    
    ax.set_xlabel('Expression régulière', fontweight='bold')
    ax.set_ylabel('Temps (ms)', fontweight='bold')
    ax.set_title('Performance de la recherche par RegEx\n(Moyenne ± Écart-type, n=50)',
                 fontweight='bold', pad=20)
    ax.set_xticks(x)
    ax.set_xticklabels(queries, rotation=15, ha='right')
    ax.grid(axis='y', alpha=0.3)
    
    # Add value labels
    for i, (bar, mean, std) in enumerate(zip(bars, means, stddevs)):
        height = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2., height + std + 1,
                f'{mean:.1f}±{std:.1f}',
                ha='center', va='bottom', fontsize=9, fontweight='bold')
    
    plt.tight_layout()
    plt.savefig('chart_regex_search.pdf', dpi=300, bbox_inches='tight')
    plt.savefig('chart_regex_search.png', dpi=300, bbox_inches='tight')
    print("✓ Generated: chart_regex_search.pdf/png")
    plt.close()

def plot_comparison(simple_data, regex_data):
    """Comparison chart between simple and regex search"""
    # Take representative queries
    simple_mean = np.mean([item['mean_ms'] for item in simple_data])
    regex_mean = np.mean([item['mean_ms'] for item in regex_data])
    
    simple_std = np.mean([item['stddev_ms'] for item in simple_data])
    regex_std = np.mean([item['stddev_ms'] for item in regex_data])
    
    fig, ax = plt.subplots(figsize=(8, 6))
    
    types = ['Recherche\nSimple', 'Recherche\nRegEx']
    means = [simple_mean, regex_mean]
    stddevs = [simple_std, regex_std]
    colors = ['#2ecc71', '#3498db']
    
    x = np.arange(len(types))
    bars = ax.bar(x, means, yerr=stddevs, capsize=8,
                   color=colors, edgecolor='black', linewidth=1.5,
                   error_kw={'linewidth': 2.5, 'ecolor': '#e74c3c'})
    
    ax.set_ylabel('Temps moyen (ms)', fontweight='bold')
    ax.set_title('Comparaison des temps de recherche\n(Moyenne globale ± Écart-type)',
                 fontweight='bold', pad=20)
    ax.set_xticks(x)
    ax.set_xticklabels(types)
    ax.grid(axis='y', alpha=0.3)
    
    # Add value labels
    for bar, mean, std in zip(bars, means, stddevs):
        height = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2., height + std + 1,
                f'{mean:.1f}±{std:.1f} ms',
                ha='center', va='bottom', fontsize=11, fontweight='bold')
    
    # Add speedup annotation
    speedup = regex_mean / simple_mean
    ax.text(0.5, max(means) * 0.8, 
            f'RegEx est {speedup:.1f}× plus lent',
            ha='center', fontsize=12, 
            bbox=dict(boxstyle='round', facecolor='yellow', alpha=0.7))
    
    plt.tight_layout()
    plt.savefig('chart_comparison.pdf', dpi=300, bbox_inches='tight')
    plt.savefig('chart_comparison.png', dpi=300, bbox_inches='tight')
    print("✓ Generated: chart_comparison.pdf/png")
    plt.close()

def plot_recommendations(data):
    """Distribution histogram of recommendation times"""
    all_times = []
    for item in data[:10]:  # Take first 10 books
        all_times.extend(item['times_ms'])
    
    fig, ax = plt.subplots(figsize=(10, 6))
    
    n, bins, patches = ax.hist(all_times, bins=30, color='#9b59b6', 
                                edgecolor='black', alpha=0.7)
    
    # Color gradient
    cm = plt.cm.get_cmap('viridis')
    bin_centers = 0.5 * (bins[:-1] + bins[1:])
    col = bin_centers - min(bin_centers)
    col /= max(col)
    for c, p in zip(col, patches):
        plt.setp(p, 'facecolor', cm(c))
    
    mean_time = np.mean(all_times)
    ax.axvline(mean_time, color='red', linestyle='--', linewidth=2.5,
               label=f'Moyenne: {mean_time:.2f} ms')
    
    ax.set_xlabel('Temps de recommandation (ms)', fontweight='bold')
    ax.set_ylabel('Fréquence', fontweight='bold')
    ax.set_title('Distribution des temps de recommandation\n(1000 mesures sur 10 livres)',
                 fontweight='bold', pad=20)
    ax.legend(fontsize=11)
    ax.grid(axis='y', alpha=0.3)
    
    plt.tight_layout()
    plt.savefig('chart_recommendations.pdf', dpi=300, bbox_inches='tight')
    plt.savefig('chart_recommendations.png', dpi=300, bbox_inches='tight')
    print("✓ Generated: chart_recommendations.pdf/png")
    plt.close()

def plot_result_count(simple_data, regex_data):
    """Bar chart showing number of results returned"""
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))
    
    # Simple search
    queries = [item['query'] for item in simple_data]
    counts = [item['result_count'] for item in simple_data]
    
    bars1 = ax1.bar(range(len(queries)), counts, color='#2ecc71',
                    edgecolor='black', linewidth=1.2)
    ax1.set_xlabel('Requête', fontweight='bold')
    ax1.set_ylabel('Nombre de livres trouvés', fontweight='bold')
    ax1.set_title('Résultats - Recherche Simple', fontweight='bold')
    ax1.set_xticks(range(len(queries)))
    ax1.set_xticklabels(queries, rotation=0)
    ax1.grid(axis='y', alpha=0.3)
    
    for bar, count in zip(bars1, counts):
        height = bar.get_height()
        ax1.text(bar.get_x() + bar.get_width()/2., height + 5,
                f'{count}', ha='center', va='bottom', fontweight='bold')
    
    # Regex search
    queries2 = [item['query'] for item in regex_data]
    counts2 = [item['result_count'] for item in regex_data]
    
    bars2 = ax2.bar(range(len(queries2)), counts2, color='#3498db',
                    edgecolor='black', linewidth=1.2)
    ax2.set_xlabel('Expression régulière', fontweight='bold')
    ax2.set_ylabel('Nombre de livres trouvés', fontweight='bold')
    ax2.set_title('Résultats - Recherche RegEx', fontweight='bold')
    ax2.set_xticks(range(len(queries2)))
    ax2.set_xticklabels(queries2, rotation=15, ha='right')
    ax2.grid(axis='y', alpha=0.3)
    
    for bar, count in zip(bars2, counts2):
        height = bar.get_height()
        ax2.text(bar.get_x() + bar.get_width()/2., height + 5,
                f'{count}', ha='center', va='bottom', fontweight='bold')
    
    plt.tight_layout()
    plt.savefig('chart_result_counts.pdf', dpi=300, bbox_inches='tight')
    plt.savefig('chart_result_counts.png', dpi=300, bbox_inches='tight')
    print("✓ Generated: chart_result_counts.pdf/png")
    plt.close()

def generate_summary_table(results):
    """Generate a summary statistics table"""
    print("\n=== SUMMARY STATISTICS ===\n")
    
    print("SIMPLE SEARCH:")
    print(f"{'Query':<15} {'Mean (ms)':<12} {'StdDev':<12} {'Min':<12} {'Max':<12} {'Results':<10}")
    print("-" * 80)
    for item in results['search_simple']:
        times = item['times_ms']
        print(f"{item['query']:<15} {item['mean_ms']:>10.2f}  {item['stddev_ms']:>10.2f}  "
              f"{min(times):>10.2f}  {max(times):>10.2f}  {item['result_count']:>8}")
    
    print("\nREGEX SEARCH:")
    print(f"{'Query':<20} {'Mean (ms)':<12} {'StdDev':<12} {'Min':<12} {'Max':<12} {'Results':<10}")
    print("-" * 85)
    for item in results['search_regex']:
        times = item['times_ms']
        print(f"{item['query']:<20} {item['mean_ms']:>10.2f}  {item['stddev_ms']:>10.2f}  "
              f"{min(times):>10.2f}  {max(times):>10.2f}  {item['result_count']:>8}")
    
    print("\nRECOMMENDATIONS:")
    times_all = [item['mean_ms'] for item in results['recommendations']]
    print(f"Average time: {np.mean(times_all):.2f} ms")
    print(f"StdDev: {np.std(times_all):.2f} ms")
    print(f"Min: {min(times_all):.2f} ms")
    print(f"Max: {max(times_all):.2f} ms")

def main():
    print("=== DAAR Project 3 - Chart Generator ===\n")
    
    # Load results
    results = load_results()
    if not results:
        return
    
    print("Generating charts...\n")
    
    # Generate all charts
    plot_simple_search(results['search_simple'])
    plot_regex_search(results['search_regex'])
    plot_comparison(results['search_simple'], results['search_regex'])
    plot_recommendations(results['recommendations'])
    plot_result_count(results['search_simple'], results['search_regex'])
    
    # Print summary
    generate_summary_table(results)
    
    print("\n✓ All charts generated successfully!")
    print("\nGenerated files:")
    print("  - chart_simple_search.pdf/png")
    print("  - chart_regex_search.pdf/png")
    print("  - chart_comparison.pdf/png")
    print("  - chart_recommendations.pdf/png")
    print("  - chart_result_counts.pdf/png")
    print("\nCopy PDF files to your LaTeX project and include them with:")
    print("  \\includegraphics[width=\\textwidth]{chart_simple_search.pdf}")

if __name__ == "__main__":
    main()