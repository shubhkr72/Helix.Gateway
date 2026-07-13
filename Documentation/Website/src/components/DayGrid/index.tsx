import React from 'react';
import Link from '@docusaurus/Link';
import styles from './styles.module.css';

export default function DayGrid() {
  const days = Array.from({ length: 29 }, (_, i) => i + 1);

  return (
    <section className={styles.sectionWrapper}>
      <div className="container">
        <div className={styles.gridContainer}>
          {days.map((day) => (
            <Link
              key={day}
              to={`/docs/list-days/day-${day}`}
              className={styles.card}
            >
              Day {day}
            </Link>
          ))}
        </div>
      </div>
    </section>
  );
}
