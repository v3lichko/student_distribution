import csv
from pathlib import Path

import numpy as np


OUTPUT_PATH = Path("raw_data/students.csv")
STUDENT_COUNT = 10_000

rng = np.random.default_rng(42)

FIRST_NAMES = [
    "Ivan", "Petr", "Alexey", "Dmitry", "Sergey",
    "Andrey", "Mikhail", "Nikita", "Artem", "Egor",
    "Anna", "Maria", "Daria", "Sofia", "Elena",
    "Polina", "Ksenia", "Alina", "Victoria", "Anastasia",
]

LAST_NAMES = [
    "Ivanov", "Petrov", "Sidorov", "Smirnov", "Kuznetsov",
    "Popov", "Vasiliev", "Sokolov", "Mikhailov", "Novikov",
    "Fedorov", "Morozov", "Volkov", "Alekseev", "Lebedev",
]

PATRONYMICS = [
    "Ivanovich", "Petrovich", "Sergeevich", "Dmitrievich",
    "Alexandrovich", "Andreevich", "Mikhailovich",
    "Ivanovna", "Petrovna", "Sergeevna", "Dmitrievna",
    "Alexandrovna", "Andreevna", "Mikhailovna",
]


def generate_full_name() -> str:
    last_name = rng.choice(LAST_NAMES)
    first_name = rng.choice(FIRST_NAMES)
    patronymic = rng.choice(PATRONYMICS)

    return f"{last_name} {first_name} {patronymic}"


def main() -> None:
    OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)

    with OUTPUT_PATH.open("w", newline="", encoding="utf-8") as file:
        writer = csv.writer(file)

        writer.writerow(["isu", "full_name", "telegram", "score"])

        for i in range(1, STUDENT_COUNT + 1):
            isu = 100000 + i
            full_name = generate_full_name()
            telegram = f"@student_{isu}"
            score = int(rng.integers(0, 101))

            writer.writerow([isu, full_name, telegram, score])

    print(f"Generated {STUDENT_COUNT} students")
    print(f"Saved to {OUTPUT_PATH}")


if __name__ == "__main__":
    main()